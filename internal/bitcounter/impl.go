package bitcounter

type Input struct {
	Data []byte
}

type Result struct {
	ByteStatistic []ByteStat `json:"byteStatistic"`
}

type ByteStat struct {
	RaisedBitQuantity int               `json:"raisedBitQuantity"`
	ByteCount         int               `json:"byteCount"`
	Values            map[int][]int     `json:"values"`
}

func countSetBits(n byte) int {
	count := 0
	for n > 0 {
		count += int(n & 1)
		n >>= 1
	}
	return count
}

func Process(in *Input) *Result {
	result := &Result{
		ByteStatistic: make([]ByteStat, 8), // Максимум 8 бит
	}

	// Инициализация структуры для хранения значений
	for i := 0; i < len(result.ByteStatistic); i++ {
		result.ByteStatistic[i].RaisedBitQuantity = i
		result.ByteStatistic[i].Values = make(map[int][]int)
	}

	// Подсчет бит и заполнение статистики
	for offset, b := range in.Data {
		raisedBitQuantity := countSetBits(b)
		result.ByteStatistic[raisedBitQuantity].ByteCount++
		if _, exists := result.ByteStatistic[raisedBitQuantity].Values[int(b)]; !exists {
			result.ByteStatistic[raisedBitQuantity].Values[int(b)] = []int{}
		}
		result.ByteStatistic[raisedBitQuantity].Values[int(b)] = append(result.ByteStatistic[raisedBitQuantity].Values[int(b)], offset)
	}

	return result
}