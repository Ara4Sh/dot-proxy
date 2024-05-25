package dnsproxy

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
)

// DNSProxy is a struct that represents a DNS proxy server.
type DNSProxy struct {
	cloudFlareDotEndpoint string
}

// NewDNSProxy is a constructor function that creates a new DNSProxy instance.
func NewDNSProxy(cloudFlareDotEndpoint string) *DNSProxy {
	return &DNSProxy{
		cloudFlareDotEndpoint: cloudFlareDotEndpoint,
	}
}

// prepareQuery is a private function that transforms a UDP DNS query to a standard TCP DNS query.
func (p *DNSProxy) prepareQuery(buffer []byte) []byte {
	log.Debug("Received UDP DNS query for prepration: ", buffer)
	msize := make([]byte, 2)
	binary.BigEndian.PutUint16(msize, uint16(len(buffer)))
	query := append(msize, buffer...)
	log.Debug("Prepared DNS query: ", query)
	return query
}

// ListenUDP is a function that listens for UDP DNS queries and forwards them to CloudFlare.
func (p *DNSProxy) ListenUDP(ctx context.Context, host string, port int) error {
	udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		log.Error("Error in resolving UDP address: ", err)
		return err
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error("Error in listening on UDP: ", err)
		return err
	}
	defer udpConn.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			udpConn.Close()
			return
		}
	}()
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return nil
		default:
			buffer := make([]byte, 1024)
			n, addr, err := udpConn.ReadFromUDP(buffer)
			log.Info("Accepting UDP client: ", addr)
			if err != nil {
				log.Error("Error in reading UDP data: ", err)
			}

			buffer = p.prepareQuery(buffer[:n])
			go func(buffer []byte, addr *net.UDPAddr) {
				response, err := p.forwardToCloudFlareDoT(buffer)

				if err != nil {
					log.Error("Error in receiving data from CloudFlare: ", err)
					return
				}

				_, err = udpConn.WriteToUDP([]byte(response[2:]), addr)
				if err != nil {
					log.Error("Error in sending data to client: ", err)
				}
			}(buffer, addr)
		}
	}

}

// ListenTCP is a function that listens for TCP DNS queries and forwards them to CloudFlare.
func (p *DNSProxy) ListenTCP(ctx context.Context, host string, port int) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, strconv.Itoa(port)))

	if err != nil {
		log.Error("Error in resolving TCP address: ", err)
		return err
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Error("Error in listening on TCP: ", err)
		return err
	}
	defer tcpListener.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			tcpListener.Close()
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return nil
		default:
			tcpConn, err := tcpListener.Accept()
			log.Info("Accepting TCP client: ", tcpConn.RemoteAddr())
			if err != nil {
				log.Error("Error in accepting tcp connection: ", err)
			}

			buffer := make([]byte, 4096)
			_, err = tcpConn.Read(buffer)
			if err != nil {
				log.Error("Error in reading TCP data: ", err)
			}

			go func(buffer []byte, conn net.Conn) {
				defer conn.Close()
				response, err := p.forwardToCloudFlareDoT(buffer)
				if err != nil {
					log.Error("Error in receiving data from CloudFlare: ", err)
					return
				}

				_, err = conn.Write(response)
				if err != nil {
					log.Error("Error in sending data to client: ", err)
				}
			}(buffer, tcpConn)
		}
	}
}

// forwardToCloudFlareDoT is a function that forwards a DNS query to CloudFlare and returns the response.
func (p *DNSProxy) forwardToCloudFlareDoT(buffer []byte) ([]byte, error) {

	conn, err := tls.Dial("tcp", p.cloudFlareDotEndpoint, nil)

	if err != nil {
		log.Error("Error in connecting to CloudFlare: ", err)
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(buffer)

	if err != nil {
		log.Error("Error in sending query to CloudFlare: ", err)
		return nil, err
	}

	response := make([]byte, 4096)
	n, err := conn.Read(response)

	if err != nil {
		log.Error("Error in reading response from CloudFlare: ", err)
		return nil, err
	}

	return response[:n], err
}
