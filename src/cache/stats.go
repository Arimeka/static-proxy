package cache

import (
	"constant"

	"github.com/jinzhu/gorm"

	"log"
	"sync"
	"time"
)

func NewStats(db *gorm.DB) (*Stats, error) {
	var (
		count uint64
		sum   float64
	)
	if err := db.Model(&File{}).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		if err := db.Raw("SELECT SUM(size) FROM files;").Row().Scan(&sum); err != nil {
			return nil, err
		}
	}

	return &Stats{
		Files:    count,
		FullSize: constant.ByteSize(sum),
		Mutex:    &sync.RWMutex{},
		DB:       db,
	}, nil
}

type Stats struct {
	Files     uint64
	FullSize  constant.ByteSize
	LimitSize constant.ByteSize

	CleanDuration time.Duration
	CleaningNow   bool

	Mutex *sync.RWMutex
	DB    *gorm.DB
}

func (s Stats) NeedCleaning() bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return !s.CleaningNow && s.FullSize > s.LimitSize
}

func (s *Stats) SetCleaningNow() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.CleaningNow = true
}

func (s *Stats) UnsetCleaningNow() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.CleaningNow = false
}

func (s *Stats) MarkToDelete() error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	var files = []*File{}

	err := s.DB.Where("deleted <> ?", true).Order("updated_at ASC").Limit(1000).Find(&files).Error
	if err != nil {
		return err
	}

	for _, file := range files {
		if s.FullSize < s.LimitSize {
			return nil
		}
		if err := s.DB.Model(file).UpdateColumn("deleted", true).Error; err != nil {
			return err
		}
		s.Files--
		s.FullSize -= constant.ByteSize(file.Size)
	}

	return nil
}

func (s Stats) DeleteFiles() error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	var files = []*File{}

	err := s.DB.Where("deleted = ?", true).Limit(1000).Find(&files).Error
	if err != nil {
		return err
	}

	for _, file := range files {
		file.Delete()
		if err := s.DB.Delete(file).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *Stats) CacheWatcher(timer *time.Timer) {
	defer timer.Stop()

	<-timer.C
	if s.NeedCleaning() {
		log.Println("Start cleaning cache...")
		s.SetCleaningNow()
		s.cleanWork()
		log.Println("Finish cleaning cache")
		s.UnsetCleaningNow()
	}

	newTimer := time.NewTimer(s.CleanDuration)
	log.Println("Next check cache after", s.CleanDuration)
	go s.CacheWatcher(newTimer)
}

func (s *Stats) cleanWork() {
	s.MarkToDelete()
	s.DeleteFiles()
	s.Mutex.RLock()
	if s.FullSize > s.LimitSize {
		timer := time.NewTimer(20 * time.Second)
		go func() {
			defer timer.Stop()
			<-timer.C
			s.cleanWork()
		}()
	}
	s.Mutex.RUnlock()
}
