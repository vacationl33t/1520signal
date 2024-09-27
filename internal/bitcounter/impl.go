package bitcounter

// Input структура для входных данных
type Input struct {
	Data []byte
}

// Result структура для хранения статистики
type Result struct {
	ByteStatistic []ByteStat `json:"byteStatistic"`
}

// ByteStat структура для статистики по байтам
type ByteStat struct {
	RaisedBitQuantity int            `json:"raisedBitQuantity"`
	ByteCount         int            `json:"byteCount"`
	Values            map[byte][]int `json:"values"`
}

// countSetBits функция для подсчета единичных бит в байте
func countSetBits(n byte) int {
	count := 0
	for n > 0 {
		count += int(n & 1)
		n >>= 1
	}
	return count
}

// Process функция для обработки входных данных и подсчета статистики
func Process(in *Input) *Result {
	result := &Result{
		ByteStatistic: make([]ByteStat, 9),
	}

	// Инициализация структуры для хранения значений
	for i := 0; i < len(result.ByteStatistic); i++ {
		result.ByteStatistic[i].RaisedBitQuantity = i
		result.ByteStatistic[i].Values = make(map[byte][]int)
	}

	// Подсчет бит и заполнение статистики
	for offset, b := range in.Data {
		raisedBitQuantity := countSetBits(b)

		// Увеличиваем количество байтов с данным raisedBitQuantity
		result.ByteStatistic[raisedBitQuantity].ByteCount++
		// Добавляем смещение в список значений для этого байта
		result.ByteStatistic[raisedBitQuantity].Values[b] = append(result.ByteStatistic[raisedBitQuantity].Values[b], offset)
	}

	return result
}
