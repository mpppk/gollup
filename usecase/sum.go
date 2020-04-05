package usecase

import (
	"log"
	"time"

	"github.com/mpppk/cli-template/domain/model"
	"github.com/mpppk/cli-template/domain/repository"
)

// Sum represents usecases related sum calculation
type Sum struct {
	sumHistoryRepository repository.SumHistory
}

// NewSum create use case related sum
func NewSum(sumHistoryRepository repository.SumHistory) *Sum {
	return &Sum{
		sumHistoryRepository: sumHistoryRepository,
	}
}

// CalcSum is use case to calculate sum
func (s *Sum) CalcSum(numbers []int) int {
	result := model.NewNumbers(numbers).CalcSum()
	now := time.Now()
	log.Printf("start saving history of sum result. date=%v, numbers=%d, result=%v\n", now, numbers, result)
	s.sumHistoryRepository.Add(&model.SumHistory{
		Date:    time.Now(), // FIXME
		Numbers: numbers,
		Result:  result,
	})
	log.Printf("finish saving history of sum result. date=%v, numbers=%d, result=%v\n", now, numbers, result)
	return result
}

// CalcL1Norm is use case to calculate L1 norm
func (s *Sum) CalcL1Norm(numbers []int) int {
	result := model.NewNumbers(numbers).CalcL1Norm()
	now := time.Now()
	log.Printf("start saving history of norm result. date=%v, numbers=%d, result=%v\n", now, numbers, result)
	s.sumHistoryRepository.Add(&model.SumHistory{
		IsNorm:  true,
		Date:    now, // FIXME
		Numbers: numbers,
		Result:  result,
	})
	log.Printf("finish saving history of norm result. date=%v, numbers=%d, result=%v\n", now, numbers, result)
	return result
}

// ListSumHistory lists history of sum
func (s *Sum) ListSumHistory(limit int) []*model.SumHistory {
	return s.sumHistoryRepository.List(limit)
}
