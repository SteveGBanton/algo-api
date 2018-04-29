package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type TwoDNA struct {
	Title1    string `json:"title1"`
	Residues1   string `json:"residues1"`
	Title2    string `json:"title2"`
	Residues2   string `json:"residues2"`
}

var sample = map[string]TwoDNA {
	"a": TwoDNA{
		Title1: "Human insulin sequence",
		Residues1: "AGCCCTCCAGGACAGGCTGCATCAGAAGAGGCCATCAAGCAGGTCTGTTCCAAGGGCCTTTGCGTCAGGTGGGCTCAGGATTCCAGGGTGGCTGGACCCCAGGCCCCAGCTCTGCAGCAGGGAGGACGTGGCTGGGCTCGTGAAGCATGTGGGGGTGAGCCCAGGGGCCCCAAGGCAGGGCACCTGGCCTTCAGCCTGCCTCAGCCCTGCCTGTCTCCCAGATCACTGTCCTTCTGCCATGGCCCTGTGGATGCGCCTCCTGCCCCTGCTGGCGCTGCTGGCCCTCTGGGGACCTGACCCAGCCGCAGCCTTTGTGAACCAACACCTGTGCGGCTCACACCTGGTGGAAGCTCTCTACCTAGTGTGCGGGGAACGAGGCTTCTTCTACACACCCAAGACCCGCCGGGAGGCAGAGGACCTGCAGGGTGAGCCAACTGCCCATTGCTGCCCCTGGCCGCCCCCAGCCACCCCCTGCTCCTGGCGCTCCCACCCAGCATGGGCAGAAGGGGGCAGGAGGCTGCCACCCAGCAGGGGGTCAGGTGCACTTTTTTAAAAAGAAGTTCTCTTGGTCACGTCCTAAAAGTGACCAGCTCCCTGTGGCCCAGTCAGAATCTCAGCCTGAGGACGGTGTTGGCTTCGGCAGCCCCGAGATACATCAGAGGGTGGGCACGCTCCTCCCTCCACTCGCCCCTCAAACAAATGCCCCGCAGCCCATTTCTCCACCCTCATTTGATGACCGCAGATTCAAGTGTTTTGTTAAGTAAAGTCCTGGGTGACCTGGGGTCACAGGGTGCCCCACGCTGCCTGCCTCTGGGCGAACACCCCATCACGCCCGGAGGAGGGCGTGGCTGCCTGCCTGAGTGGGCCAGACCCCTGTCGCCAGGCCTCACGGCAGCTCCATAGTCAGGAGATGGGGAAGATGCTGGGGACAGGCCCTGGGGAGAAGTACTGGGATCACCTGTTCAGGCTCCCACTGTGACGCTGCCCCGGGGCGGGGGAAGGAGGTGGGACATGTGGGCGTTGGGGCCTGTAGGTCCACACCCAGTGTGGGTGACCCTCCCTCTAACCTGGGTCCAGCCCGGCTGGAGATGGGTGGGAGTGCGACCTAGGGCTGGCGGGCAGGCGGGCACTGTGTCTCCCTGACTGTGTCCTCCTGTGTCCCTCTGCCTCGCCGCTGTTCCGGAACCTGCTCTGCGCGGCACGTCCTGGCAGTGGGGCAGGTGGAGCTGGGCGGGGGCCCTGGTGCAGGCAGCCTGCAGCCCTTGGCCCTGGAGGGGTCCCTGCAGAAGCGTGGCATTGTGGAACAATGCTGTACCAGCATCTGCTCCCTCTACCAGCTGGAGAACTACTGCAACTAGACGCAGCCCGCAGGCAGCCCCACACCCGCCGCCTCCTGCACCGAGAGAGATGGAATAAAGCCCTTGAACCAGC",
		Title2: "Mouse insulin sequence",
		Residues2: "ACCAGGCAAGTGTTTGGAAACTGCAGCTTCAGCCCCTCTGGCCATCTGCCTACCCACCCCACCTGGAGACCTTAATGGGCCAAACAGCAAAGTCCAGGGGGCAGAGAGGAGGTACTTTGGACTATAAAGCTGGTGGGCATCCAGTAACCCCCAGCCCTTAGTGACCAGCTATAATCAGAGACCATCAGCAAGCAGGTATGTACTCTCCTCTTTGGGCCTGGCTCCCCAGCCAAGACTCCAGCGACTTTAGGGAGAATGTGGGCTCCTCTCTTACATGGATCTTTTGCTAGCCTCAACCCTGCCTATCTTTCAGGTCATTGTTTCAACATGGCCCTGTTGGTGCACTTCCTACCCCTGCTGGCCCTGCTTGCCCTCTGGGAGCCCAAACCCACCCAGGCTTTTGTCAAACAGCATCTTTGTGGTCCCCACCTGGTAGAGGCTCTCTACCTGGTGTGTGGGGAGCGTGGCTTCTTCTACACACCCAAGTCCCGCCGTGAAGTGGAGGACCCACAAGTGGAACAACTGGAGCTGGGAGGAAGCCCCGGGGACCTTCAGACCTTGGCGTTGGAGGTGGCCCGGCAGAAGCGTGGCATTGTGGATCAGTGCTGCACCAGCATCTGCTCCCTCTACCAGCTGGAGAACTACTGCAACTAAGGCCCACCTCGACCCGCCCCACCCCTCTGCAATGAATAAAACTTTTGAATAAGCACCAAAAAAAA",
	},
}

func DNACompare(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
		case http.MethodGet:
			stringA := sample["a"].Residues1
			stringB := sample["a"].Residues2
			sequenceLength, tableA := DynamicLCS(stringA, stringB)
			sequenceA, sequenceB, longestCS := GetSequence(tableA, stringA, stringB)
			writeJSON(w, sequenceLength)
			writeJSON(w, sequenceA)
			writeJSON(w, sequenceB)
			writeJSON(w, longestCS)
			writeJSON(w, tableA)
		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			sequences := FromJSON(body)
			stringA := sequences.Residues1
			stringB := sequences.Residues2
			sequenceLength, tableA := DynamicLCS(stringA, stringB)
			GetSequence(tableA, stringA, stringB)
			writeJSON(w, sequenceLength)
			//writeJSON(w, sequenceA)
			// writeJSON(w, sequenceB)
			//writeJSON(w, longestCS)
			// writeJSON(w, tableA)
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unsupported request method."))
	}
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func ToJSON(i interface{}) []byte {
	ToJSON, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return ToJSON
}

func FromJSON(data []byte) TwoDNA {
	book := TwoDNA{}
	err := json.Unmarshal(data, &book)
	if err != nil {
		panic(err)
	}
	return book
}

func DynamicLCS(A string, B string) (int, [][]int) {
	lengthA := len(A)
	lengthB := len(B)

	// Make a 2D tables of height lengthA, width lengthB
	tableA := make([][]int, lengthA + 1)
	tableB := make([][]int, lengthA + 1)

	for i := 0; i < (lengthA + 1); i++ {
		tableA[i] = make([]int, lengthB + 1)
		tableB[i] = make([]int, lengthB + 1)
	}

	var runesOfA = []rune(A)
	var runesOfB = []rune(B)

	for i := 1; i <= lengthA; i++ {
		for j := 1; j <= lengthB; j++ {
			if (string(runesOfA[i - 1]) == string(runesOfB[j - 1])) {
				tableB[i][j] = tableB[i - 1][j - 1] + 1;
				tableA[i][j] = 2;
			} else if (tableB[i - 1][j] >= tableB[i][j - 1]) {
				tableB[i][j] = tableB[i - 1][j];
				tableA[i][j] = 3;
			} else {
				tableB[i][j] = tableB[i][j - 1];
				tableA[i][j] = 1;
			}
		}
	}

	sequenceLength := tableB[lengthA - 1][lengthB - 1]

	return sequenceLength, tableA
}

func GetSequence(tableA [][]int, A string, B string) ([]string, []string, []string) {
	sequenceA := []string{}
	sequenceB := []string{}
	longestCS := []string{}
	bCount := len(B)
	aCount := len(A)

	for aCount > 0 && bCount > 0 {
		if (tableA[aCount][bCount] == 2) {
			aCount -= 1
			bCount -= 1
			char := string(A[aCount])
			sequenceA = append([]string{char}, sequenceA...)
			sequenceB = append([]string{char}, sequenceB...)
			longestCS = append([]string{char}, longestCS...)
		} else if (tableA[aCount][bCount] == 3) {
			aCount -= 1;
			sequenceA = append([]string{"*"}, sequenceA...)
		} else {
			bCount -= 1;
			sequenceB = append([]string{"*"}, sequenceB...)
		}
	}

	for len(sequenceA) < len(A) {
		sequenceA = append([]string{"*"}, sequenceA...)
	}
	for len(sequenceB) < len(B) {
		sequenceB = append([]string{"*"}, sequenceB...)
	}
	return sequenceA, sequenceB, longestCS
}
