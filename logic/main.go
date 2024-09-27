package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("1. Soal no 1")
	fmt.Println("2. Soal no 2")
	fmt.Println("3. Soal no 3")
	fmt.Println("4. Soal no 4")
	fmt.Print("Masukkan pilihan: ")

	var pil int
	fmt.Scanln(&pil)

	switch pil {
	case 1:
		soal1()
	case 2:
		soal2()
	case 3:
		soal3()
	case 4:
		soal4()
	default:
		fmt.Println("Pilihan tidak valid.")
	}
}

func soal1() {
	var N int
	fmt.Print("Jumlah string: ")
	fmt.Scanln(&N)

	stringsList := make([]string, N)
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < N; i++ {
		input, _ := reader.ReadString('\n')
		stringsList[i] = strings.TrimSpace(strings.ToLower(input))
	}

	if result1, result2, found := findMatchingStrings(stringsList); found {
		fmt.Printf("Result: %d %d\n", result1, result2)
	} else {
		fmt.Println("Result: false")
	}
}

func findMatchingStrings(stringsList []string) (int, int, bool) {
	stringSet := make(map[string]int)

	for i, str := range stringsList {
		if idx, exists := stringSet[str]; exists {
			return idx + 1, i + 1, true
		}
		stringSet[str] = i
	}
	return 0, 0, false
}

func soal2() {
	var totalBelanja, jumlahDibayarkan int
	fmt.Print("total belanja: ")
	fmt.Scanln(&totalBelanja)
	fmt.Print("jumlah dibayarkan: ")
	fmt.Scanln(&jumlahDibayarkan)

	kembalianAsli, kembalianBulat, pecahan := hitungKembalian(totalBelanja, jumlahDibayarkan)

	if kembalianAsli == -1 {
		fmt.Println("False, kurang bayar")
	} else {
		fmt.Printf("Kembalian yang harus diberikan kasir: %d, dibulatkan menjadi %d\n", kembalianAsli, kembalianBulat)
		fmt.Println("Pecahan uang:")
		for nilai, jumlah := range pecahan {
			if jumlah > 0 {
				if nilai >= 1000 {
					fmt.Printf("%d lembar %d\n", jumlah, nilai)
				} else {
					fmt.Printf("%d koin %d\n", jumlah, nilai)
				}
			}
		}
	}
}

func hitungKembalian(totalBelanja, jumlahDibayarkan int) (int, int, map[int]int) {
	pecahanUang := []int{100000, 50000, 20000, 10000, 5000, 2000, 1000, 500, 200, 100}
	pecahanDetail := make(map[int]int)

	kembalianAsli := jumlahDibayarkan - totalBelanja

	if kembalianAsli < 0 {
		return -1, 0, nil
	}

	kembalianBulat := (kembalianAsli / 100) * 100

	for _, pecahan := range pecahanUang {
		if kembalianBulat >= pecahan {
			pecahanDetail[pecahan] = kembalianBulat / pecahan
			kembalianBulat %= pecahan
		}
	}

	return kembalianAsli, (jumlahDibayarkan - totalBelanja) / 100 * 100, pecahanDetail
}

func soal3() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Masukkan string yang ingin divalidasi (ketik 'exit' untuk berhenti):")

	for {
		fmt.Print("Input: ")
		scanner.Scan()
		input := scanner.Text()

		if strings.ToLower(input) == "exit" {
			break
		}

		if validateString(input) {
			fmt.Println("Valid: true")
		} else {
			fmt.Println("Valid: false")
		}
	}
}

func validateString(input string) bool {
	if len(input) < 1 || len(input) > 4096 {
		return false
	}

	pairs := map[rune]rune{
		'}': '{',
		']': '[',
		'>': '<',
	}

	var stack []rune
	for _, char := range input {
		if isOpeningBracket(char) {
			stack = append(stack, char)
		} else if isClosingBracket(char) {
			if len(stack) == 0 || stack[len(stack)-1] != pairs[char] {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			return false
		}
	}
	return len(stack) == 0
}

func isOpeningBracket(char rune) bool {
	return char == '{' || char == '[' || char == '<'
}

func isClosingBracket(char rune) bool {
	return char == '}' || char == ']' || char == '>'
}

func soal4() {
	var collectiveLeave int
	var joinDateStr, plannedLeaveStr string
	var leaveDuration int

	fmt.Print("Masukkan jumlah cuti bersama: ")
	fmt.Scan(&collectiveLeave)

	fmt.Print("Masukkan tanggal join karyawan (YYYY-MM-DD): ")
	fmt.Scan(&joinDateStr)

	fmt.Print("Masukkan tanggal rencana cuti (YYYY-MM-DD): ")
	fmt.Scan(&plannedLeaveStr)

	fmt.Print("Masukkan durasi cuti (hari): ")
	fmt.Scan(&leaveDuration)

	joinDate, err1 := time.Parse("2006-01-02", joinDateStr)
	plannedLeaveDate, err2 := time.Parse("2006-01-02", plannedLeaveStr)

	if err1 != nil || err2 != nil {
		fmt.Println("Format tanggal tidak valid. Pastikan menggunakan format YYYY-MM-DD.")
		return
	}

	canTake, reason := canTakePersonalLeave(joinDate, collectiveLeave, plannedLeaveDate, leaveDuration)

	if canTake {
		fmt.Println("True")
	} else {
		fmt.Println("False")
		fmt.Println("Alasan:", reason)
	}
}

func canTakePersonalLeave(joinDate time.Time, collectiveLeave int, plannedLeaveDate time.Time, leaveDuration int) (bool, string) {
	firstEligibleDate := joinDate.Add(180 * 24 * time.Hour)

	if plannedLeaveDate.Before(firstEligibleDate) {
		return false, "Karena belum 180 hari sejak tanggal join karyawan"
	}

	endOfYear := time.Date(joinDate.Year(), time.December, 31, 0, 0, 0, 0, joinDate.Location())
	daysEligible := int(endOfYear.Sub(firstEligibleDate).Hours() / 24)

	totalLeave := (daysEligible * collectiveLeave) / 365

	if leaveDuration > totalLeave {
		return false, fmt.Sprintf("Karena hanya boleh mengambil %d hari cuti", totalLeave)
	}
	if leaveDuration > 3 {
		return false, "Karena maksimal cuti pribadi adalah 3 hari berturut-turut"
	}

	return true, ""
}
