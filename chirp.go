package main

import "strings"

func FilterChirp(c string) string {

	bannedWords := make(map[string]interface{})
	bannedWords["kerfuffle"] = nil
	bannedWords["sharbert"] = nil
	bannedWords["fornax"] = nil

	chirpWords := strings.Split(c, " ")

	var checkWord string
	var ok bool

	for w := 0; w < len(chirpWords); w++ {
		checkWord = strings.ToLower(chirpWords[w])
		if _, ok = bannedWords[checkWord]; ok {
			chirpWords[w] = "****"
		}
	}

	return strings.Join(chirpWords, " ")

}
