package storage

//import (
//	"bufio"
//	"encoding/json"
//	"errors"
//	"os"
//	"sync"
//
//	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/configs"
//	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/app/helpers"
//)
//
//type Producer interface {
//	WriteEvent(event *storage)
//	Close() error
//}
//
//type Consumer interface {
//	ReadEvent() (*storage, error)
//	Close() error
//}
//
//type storage struct {
//	Data map[UserID]ShortLinks
//	mu   sync.Mutex
//}
//
//type producer struct {
//	file    *os.File
//	write   *bufio.Writer
//	encoder *json.Encoder
//}
//
//type consumer struct {
//	file    *os.File
//	read    *bufio.Reader
//	decoder *json.Decoder
//}
//
//func (s *storage) LinkBy(userID UserID, sl ShortLink) (string, error) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	shortLinks, ok := s.Data[userID]
//	if !ok {
//		return "", errors.New("url not found")
//	}
//
//	url, ok := shortLinks[sl]
//	if !ok {
//		return "", errors.New("url not found")
//	}
//
//	return url, nil
//}
//
//func (s *storage) LinksByUser(userID UserID) (ShortLinks, error) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	shortLinks, ok := s.Data[userID]
//
//	if !ok {
//		return shortLinks, errors.New("url not found")
//	}
//
//	return shortLinks, nil
//}
//
//func (s *storage) Save(userID UserID, url string) ShortLink {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	sl := ShortLink(helpers.RandomString(10))
//
//	currentUrls := ShortLinks{}
//
//	if urls, ok := s.Data[userID]; ok {
//		currentUrls = urls
//	}
//
//	currentUrls[sl] = url
//
//	s.Data[userID] = currentUrls
//	s.Data["default"] = currentUrls
//
//	return sl
//}
//
//func NewProducer(filename string) (*producer, error) {
//	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return &producer{
//		file:    file,
//		write:   bufio.NewWriter(file),
//		encoder: json.NewEncoder(file),
//	}, nil
//}
//
//func NewConsumer(filename string) (*consumer, error) {
//	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0777)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return &consumer{
//		file:    file,
//		read:    bufio.NewReader(file),
//		decoder: json.NewDecoder(file),
//	}, nil
//}
//
//func (p *producer) Close() error {
//	return p.file.Close()
//}
//
//func (p *consumer) Close() error {
//	return p.file.Close()
//}
//
//func (s *storage) Load(c configs.Config) error {
//	if c.FileStoragePath == "" {
//		return nil
//	}
//
//	cns, err := NewConsumer(c.FileStoragePath)
//
//	if err != nil {
//		return err
//	}
//
//	cns.decoder.Decode(&s.Data)
//
//	return nil
//}
//
//func (s *storage) Flush(c configs.Config) error {
//	if c.FileStoragePath == "" {
//		return nil
//	}
//
//	p, err := NewProducer(c.FileStoragePath)
//
//	if err != nil {
//		return err
//	}
//
//	p.encoder.Encode(&s.Data)
//
//	return p.write.Flush()
//}
//
//func New() *storage {
//	return &storage{
//		Data: make(map[UserID]ShortLinks),
//	}
//}
