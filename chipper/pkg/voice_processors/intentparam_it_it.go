package wirepod

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

func paramCheckerSlotsItIT(req interface{}, intent string, slots map[string]string, isOpus bool, justThisBotNum int, botSerial string) {
	// var req2 *vtt.IntentRequest
	// var req1 *vtt.KnowledgeGraphRequest
	// var req3 *vtt.IntentGraphRequest
	// if str, ok := req.(*vtt.IntentRequest); ok {
	// 	req2 = str
	// } else if str, ok := req.(*vtt.KnowledgeGraphRequest); ok {
	// 	req1 = str
	// } else if str, ok := req.(*vtt.IntentGraphRequest); ok {
	// 	req3 = str
	// }
	var intentParam string
	var intentParamValue string
	var newIntent string
	var isParam bool
	var intentParams map[string]string
	var botLocation string = "San Francisco"
	var botUnits string = "F"
	var botPlaySpecific bool = false
	var botIsEarlyOpus bool = false
	logger("paramCheckerSlotsItIT")

	if _, err := os.Stat("./botConfig.json"); err == nil {
		type botConfigJSON []struct {
			ESN             string `json:"ESN"`
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}
		byteValue, err := os.ReadFile("./botConfig.json")
		if err != nil {
			logger(err)
		}
		var botConfig botConfigJSON
		json.Unmarshal(byteValue, &botConfig)
		for _, bot := range botConfig {
			if strings.ToLower(bot.ESN) == botSerial {
				logger("Found bot config for " + bot.ESN)
				botLocation = bot.Location
				botUnits = bot.Units
				botPlaySpecific = bot.UsePlaySpecific
				botIsEarlyOpus = bot.IsEarlyOpus
			}
		}
	}
	if strings.Contains(intent, "volume") {
		if slots["volume"] != "" {
			newIntent = "intent_imperative_volumelevel_extend"
			isParam = true
			intentParam = "volume_level"
			if strings.Contains(slots["volume"], "medio basso") {
				intentParamValue = "VOLUME_2"
			} else if strings.Contains(slots["volume"], "basso") {
				intentParamValue = "VOLUME_1"
			} else if strings.Contains(slots["volume"], "medio alto") {
				intentParamValue = "VOLUME_4"
			} else if strings.Contains(slots["volume"], "alto") {
				intentParamValue = "VOLUME_5"
			} else if strings.Contains(slots["volume"], "medio") {
				intentParamValue = "VOLUME_3"
			} else {
				intentParamValue = "VOLUME_1"
			}
		} else {
			isParam = false
			intentParam = ""
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "eyecolor") {
		isParam = true
		newIntent = "intent_imperative_eyecolor_specific_extend"
		intentParam = "eye_color"
		if strings.Contains(slots["eye_color"], "viola") || strings.Contains(slots["eye_color"], "lilla") {
			intentParamValue = "COLOR_PURPLE"
		} else if strings.Contains(slots["eye_color"], "blu") {
			intentParamValue = "COLOR_BLUE"
		} else if strings.Contains(slots["eye_color"], "giallo") {
			intentParamValue = "COLOR_YELLOW"
		} else if strings.Contains(slots["eye_color"], "verde acqua") {
			intentParamValue = "COLOR_TEAL"
		} else if strings.Contains(slots["eye_color"], "verde") || strings.Contains(slots["eye_color"], "verdi") {
			intentParamValue = "COLOR_GREEN"
		} else if strings.Contains(slots["eye_color"], "arancione") || strings.Contains(slots["eye_color"], "arancio") {
			intentParamValue = "COLOR_ORANGE"
		} else {
			newIntent = intent
			intentParamValue = ""
			isParam = false
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "_selfie") {
		newIntent = "intent_photo_take_extend"
		intentParam = "entity_photo_selfie"
		intentParamValue = "photo_selfie"
		isParam = true
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "_noselfie") {
		newIntent = "intent_photo_take_extend"
		intentParam = "entity_photo_selfie"
		intentParamValue = ""
		isParam = true
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "settimer") {
		isParam = true
		newIntent = intent
		slotNum := slots["num"]
		slotUnit := slots["unit"]
		timerSecs, err := strconv.Atoi(slotNum)
		if err != nil {
			logger(err)
		}
		if slotNum != "" && slotUnit != "" {
			if strings.Contains(slotUnit, "minuto") || strings.Contains(slotUnit, "minuti") {
				timerSecs = timerSecs * 60
			} else if strings.Contains(slotUnit, "ora") || strings.Contains(slotUnit, "ore") {
				timerSecs = timerSecs * 60 * 60
			}
		}
		logger("Seconds parsed from speech: " + strconv.Itoa(timerSecs))
		intentParam = "timer_duration"
		intentParamValue = strconv.Itoa(timerSecs)
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "global_stop_extend") {
		isParam = true
		newIntent = intent
		intentParam = "what_to_stop"
		intentParamValue = "timer"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_knowledgegraph_prompt") {
		isParam = false
		newIntent = "intent_knowledge_promptquestion"
		intentParam = ""
		intentParamValue = ""
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_weather_extend") {
		isParam = true
		newIntent = intent
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit := weatherParser("no", botLocation, botUnits)
		intentParams = map[string]string{"condition": condition, "is_forecast": is_forecast, "local_datetime": local_datetime, "speakable_location_string": speakable_location_string, "temperature": temperature, "temperature_unit": temperature_unit}
	} else {
		if intentParam == "" {
			newIntent = intent
			intentParam = ""
			intentParamValue = ""
			isParam = false
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	if isOpus || botIsEarlyOpus || botPlaySpecific {
		if strings.Contains(intent, "intent_play_blackjack") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "blackjack"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_fistbump") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "fist_bump"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_rollcube") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "roll_cube"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_praise") {
			isParam = false
			newIntent = "intent_imperative_affirmative"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_love") {
			isParam = false
			newIntent = "intent_greeting_hello"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_abuse") {
			isParam = false
			newIntent = "intent_imperative_negative"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	IntentPass(req, newIntent, intent, intentParams, isParam, justThisBotNum)
}

func paramCheckerItIT(req interface{}, intent string, speechText string, justThisBotNum int, botSerial string) {
	var intentParam string
	var intentParamValue string
	var newIntent string
	var isParam bool
	var intentParams map[string]string
	var botLocation string = "San Francisco"
	var botUnits string = "F"
	var botPlaySpecific bool = false
	var botIsEarlyOpus bool = false
	logger("paramCheckerItIT")
	if _, err := os.Stat("./botConfig.json"); err == nil {
		type botConfigJSON []struct {
			ESN             string `json:"ESN"`
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}
		byteValue, err := os.ReadFile("./botConfig.json")
		if err != nil {
			logger(err)
		}
		var botConfig botConfigJSON
		json.Unmarshal(byteValue, &botConfig)
		for _, bot := range botConfig {
			if strings.ToLower(bot.ESN) == botSerial {
				logger("Found bot config for " + bot.ESN)
				botLocation = bot.Location
				botUnits = bot.Units
				botPlaySpecific = bot.UsePlaySpecific
				botIsEarlyOpus = bot.IsEarlyOpus
			}
		}
	}
	if botPlaySpecific {
		if strings.Contains(intent, "intent_play_blackjack") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "blackjack"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_fistbump") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "fist_bump"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_rollcube") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "roll_cube"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_popawheelie") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "pop_a_wheelie"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_pickupcube") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "pick_up_cube"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_play_keepaway") {
			isParam = true
			newIntent = "intent_play_specific_extend"
			intentParam = "entity_behavior"
			intentParamValue = "keep_away"
			intentParams = map[string]string{intentParam: intentParamValue}
		} else {
			newIntent = intent
			intentParam = ""
			intentParamValue = ""
			isParam = false
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	logger("Checking params for candidate intent " + intent)
	if strings.Contains(intent, "intent_photo_take_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, "fammi una foto") || strings.Contains(speechText, "scattami una foto") || strings.Contains(speechText, "scatta una foto") || strings.Contains(speechText, "fa una foto") {
			intentParam = "entity_photo_selfie"
			intentParamValue = "photo_selfie"
		} else {
			intentParam = "entity_photo_selfie"
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_eyecolor") {
		isParam = true
		newIntent = "intent_imperative_eyecolor_specific_extend"
		intentParam = "eye_color"
		if strings.Contains(speechText, "viola") || strings.Contains(speechText, "lilla") {
			intentParamValue = "COLOR_PURPLE"
		} else if strings.Contains(speechText, "blu") {
			intentParamValue = "COLOR_BLUE"
		} else if strings.Contains(speechText, "giallo") || strings.Contains(speechText, "gialli") {
			intentParamValue = "COLOR_YELLOW"
		} else if strings.Contains(speechText, "verde acqua") {
			intentParamValue = "COLOR_TEAL"
		} else if strings.Contains(speechText, "verde") || strings.Contains(speechText, "verdi") {
			intentParamValue = "COLOR_GREEN"
		} else if strings.Contains(speechText, "arancione") || strings.Contains(speechText, "arancioni") || strings.Contains(speechText, "arancio") {
			intentParamValue = "COLOR_ORANGE"
		} else {
			newIntent = intent
			intentParamValue = ""
			isParam = false
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_weather_extend") {
		isParam = true
		newIntent = intent
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit := weatherParser(speechText, botLocation, botUnits)
		intentParams = map[string]string{"condition": condition, "is_forecast": is_forecast, "local_datetime": local_datetime, "speakable_location_string": speakable_location_string, "temperature": temperature, "temperature_unit": temperature_unit}
	} else if strings.Contains(intent, "intent_imperative_volumelevel_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, "medio basso") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_2"
		} else if strings.Contains(speechText, "basso") || strings.Contains(speechText, "silenzioso") || strings.Contains(speechText, "minimo") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else if strings.Contains(speechText, "medio alto") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_4"
		} else if strings.Contains(speechText, "medio") || strings.Contains(speechText, "normale") || strings.Contains(speechText, "standard") || strings.Contains(speechText, "regolare") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_3"
		} else if strings.Contains(speechText, "alto") || strings.Contains(speechText, "rumoroso") || strings.Contains(speechText, "massimo") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_5"
		} else if strings.Contains(speechText, "muto") || strings.Contains(speechText, "zero") || strings.Contains(speechText, "silenzioso") || strings.Contains(speechText, "spento") {
			// there is no VOLUME_0 :(
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_names_username_extend") {
		var username string
		var nameSplitter string
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, " è ") {
			nameSplitter = " è "
		} else if strings.Contains(speechText, "sono") {
			nameSplitter = "sono"
		} else if strings.Contains(speechText, "chiamo") {
			nameSplitter = "chiamo"
		}
		if strings.Contains(speechText, " è ") || strings.Contains(speechText, "sono") || strings.Contains(speechText, "chiamo") {
			splitPhrase := strings.SplitAfter(speechText, nameSplitter)
			username = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				username = username + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			logger("Name parsed from speech: " + "`" + username + "`")
			intentParam = "username"
			intentParamValue = username
			intentParams = map[string]string{intentParam: intentParamValue}
		} else {
			logger("No name parsed from speech")
			intentParam = "username"
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	} else if strings.Contains(intent, "intent_clock_settimer_extend") {
		isParam = true
		newIntent = intent
		timerSecs := words2num(speechText)
		logger("Seconds parsed from speech: " + timerSecs)
		intentParam = "timer_duration"
		intentParamValue = timerSecs
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_global_stop_extend") {
		isParam = true
		newIntent = intent
		intentParam = "what_to_stop"
		intentParamValue = "timer"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_playmessage_extend") {
		var given_name string
		isParam = true
		newIntent = intent
		intentParam = "given_name"
		if strings.Contains(speechText, " per ") {
			splitPhrase := strings.SplitAfter(speechText, " per ")
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_recordmessage_extend") {
		var given_name string
		isParam = true
		newIntent = intent
		intentParam = "given_name"
		if strings.Contains(speechText, " per ") {
			splitPhrase := strings.SplitAfter(speechText, " per ")
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else {
		if intentParam == "" {
			newIntent = intent
			intentParam = ""
			intentParamValue = ""
			isParam = false
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	if botIsEarlyOpus {
		if strings.Contains(intent, "intent_imperative_praise") {
			isParam = false
			newIntent = "intent_imperative_affirmative"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_abuse") {
			isParam = false
			newIntent = "intent_imperative_negative"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		} else if strings.Contains(intent, "intent_imperative_love") {
			isParam = false
			newIntent = "intent_greeting_hello"
			intentParam = ""
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	}
	IntentPass(req, newIntent, speechText, intentParams, isParam, justThisBotNum)
}

func prehistoricParamCheckerItIT(req interface{}, intent string, speechText string, justThisBotNum int, botSerial string) {
	// intent.go detects if the stream uses opus or PCM.
	// If the stream is PCM, it is likely a bot with 0.10.
	// This accounts for the newer 0.10.1### builds.
	var intentParam string
	var intentParamValue string
	var newIntent string
	var isParam bool
	var intentParams map[string]string
	var botLocation string = "San Francisco"
	var botUnits string = "F"
	logger("prehistoricParamCheckerItIT")
	if _, err := os.Stat("./botConfig.json"); err == nil {
		type botConfigJSON []struct {
			ESN             string `json:"ESN"`
			Location        string `json:"location"`
			Units           string `json:"units"`
			UsePlaySpecific bool   `json:"use_play_specific"`
			IsEarlyOpus     bool   `json:"is_early_opus"`
		}
		byteValue, err := os.ReadFile("./botConfig.json")
		if err != nil {
			logger(err)
		}
		var botConfig botConfigJSON
		json.Unmarshal(byteValue, &botConfig)
		for _, bot := range botConfig {
			if strings.ToLower(bot.ESN) == botSerial {
				logger("Found bot config for " + bot.ESN)
				botLocation = bot.Location
				botUnits = bot.Units
			}
		}
	}
	if strings.Contains(intent, "intent_photo_take_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, "fammi una foto") || strings.Contains(speechText, "scattami una foto") || strings.Contains(speechText, "scatta una foto") || strings.Contains(speechText, "fa una foto") {
			intentParam = "entity_photo_selfie"
			intentParamValue = "photo_selfie"
		} else {
			intentParam = "entity_photo_selfie"
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_eyecolor") {
		// leaving stuff like this in case someone wants to add features like this to older software
		isParam = true
		newIntent = "intent_imperative_eyecolor_specific_extend"
		intentParam = "eye_color"
		if strings.Contains(speechText, "viola") || strings.Contains(speechText, "lilla") {
			intentParamValue = "COLOR_PURPLE"
		} else if strings.Contains(speechText, "blu") {
			intentParamValue = "COLOR_BLUE"
		} else if strings.Contains(speechText, "giallo") || strings.Contains(speechText, "gialli") {
			intentParamValue = "COLOR_YELLOW"
		} else if strings.Contains(speechText, "verde acqua") {
			intentParamValue = "COLOR_TEAL"
		} else if strings.Contains(speechText, "verde") || strings.Contains(speechText, "verdi") {
			intentParamValue = "COLOR_GREEN"
		} else if strings.Contains(speechText, "arancione") || strings.Contains(speechText, "arancioni") || strings.Contains(speechText, "arancio") {
			intentParamValue = "COLOR_ORANGE"
		} else {
			newIntent = intent
			intentParamValue = ""
			isParam = false
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_weather_extend") {
		isParam = true
		newIntent = intent
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit := weatherParser(speechText, botLocation, botUnits)
		intentParams = map[string]string{"condition": condition, "is_forecast": is_forecast, "local_datetime": local_datetime, "speakable_location_string": speakable_location_string, "temperature": temperature, "temperature_unit": temperature_unit}
	} else if strings.Contains(intent, "intent_imperative_volumelevel_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, "medio basso") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_2"
		} else if strings.Contains(speechText, "basso") || strings.Contains(speechText, "silenzioso") || strings.Contains(speechText, "minimo") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else if strings.Contains(speechText, "medio alto") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_4"
		} else if strings.Contains(speechText, "medio") || strings.Contains(speechText, "normale") || strings.Contains(speechText, "standard") || strings.Contains(speechText, "regolare") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_3"
		} else if strings.Contains(speechText, "alto") || strings.Contains(speechText, "rumoroso") || strings.Contains(speechText, "massimo") {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_5"
		} else if strings.Contains(speechText, "muto") || strings.Contains(speechText, "zero") || strings.Contains(speechText, "silenzioso") || strings.Contains(speechText, "spento") {
			// there is no VOLUME_0 :(
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_names_username_extend") {
		var username string
		var nameSplitter string
		isParam = true
		newIntent = "intent_names_username"
		if strings.Contains(speechText, " è ") {
			nameSplitter = " è "
		} else if strings.Contains(speechText, "sono") {
			nameSplitter = "sono"
		} else if strings.Contains(speechText, "chiamo") {
			nameSplitter = "chiamo"
		}
		if strings.Contains(speechText, " è ") || strings.Contains(speechText, "sono") || strings.Contains(speechText, "chiamo") {
			splitPhrase := strings.SplitAfter(speechText, nameSplitter)
			username = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				username = username + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			logger("Name parsed from speech: " + "`" + username + "`")
			intentParam = "username"
			intentParamValue = username
			intentParams = map[string]string{intentParam: intentParamValue}
		} else {
			logger("No name parsed from speech")
			intentParam = "username"
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	} else if strings.Contains(intent, "intent_clock_settimer_extend") {
		isParam = true
		newIntent = "intent_clock_settimer"
		timerSecs := words2num(speechText)
		logger("Seconds parsed from speech: " + timerSecs)
		intentParam = "timer_duration"
		intentParamValue = timerSecs
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_global_stop_extend") {
		isParam = true
		newIntent = "intent_global_stop"
		intentParam = "what_to_stop"
		intentParamValue = "timer"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_playmessage_extend") {
		var given_name string
		isParam = true
		newIntent = "intent_message_playmessage"
		intentParam = "given_name"
		if strings.Contains(speechText, " per ") {
			splitPhrase := strings.SplitAfter(speechText, " per ")
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_recordmessage_extend") {
		var given_name string
		isParam = true
		newIntent = "intent_message_recordmessage"
		intentParam = "given_name"
		if strings.Contains(speechText, " per ") {
			splitPhrase := strings.SplitAfter(speechText, " per ")
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_blackjack") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "blackjack"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_fistbump") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "fist_bump"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_rollcube") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "roll_cube"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_praise") {
		isParam = false
		newIntent = "intent_imperative_affirmative"
		intentParam = ""
		intentParamValue = ""
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_abuse") {
		isParam = false
		newIntent = "intent_imperative_negative"
		intentParam = ""
		intentParamValue = ""
		intentParams = map[string]string{intentParam: intentParamValue}
	} else {
		newIntent = intent
		intentParam = ""
		intentParamValue = ""
		isParam = false
		intentParams = map[string]string{intentParam: intentParamValue}
	}
	IntentPass(req, newIntent, speechText, intentParams, isParam, justThisBotNum)
}