package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type ODB2_data struct {
	Vin                                 string  `json:"vin"`
	Engine_load                         int     `json:"engine_load"`
	Engine_coolant_temperature          int     `json:"engine_coolant_temperature"`
	Fuel_pressure                       int     `json:"fuel_pressure"`
	Engine_speed                        float32 `json:"engine_speed"`
	Vehicle_speed                       int     `json:"vehicle_speed"`
	Intake_air_temperature              int     `json:"intake_air_temperature"`
	Mass_air_flow                       int     `json:"mass_air_flow"`
	Fuel_tank_level_input               int     `json:"fuel_tank_level_input"`
	Absolute_barometric_pressure        int     `json:"absolute_barometric_pressure"`
	Control_module_voltage              float32 `json:"control_module_voltage"`
	Throttle_position                   int     `json:"throttle_position"`
	Ambient_air_temperature             int     `json:"ambient_air_temperature"`
	Relative_accelerator_pedal_position int     `json:"relative_accelerator_pedal_position"`
	Engine_oil_temperature              int     `json:"engine_oil_temperature"`
	Engine_fuel_rate                    float32 `json:"engine_fuel_rate"`
	Service_distance                    int     `json:"service_distance"`
	Anti_lock_barking_active            int     `json:"anti_lock_barking_active"`
	Steering_wheel_angle                int     `json:"steering_wheel_angle"`
	Position_of_doors                   int     `json:"position_of_doors"`
	Right_left_turn_signal_light        int     `json:"right_left_turn_signal_light"`
	Alternate_beam_head_light           int     `json:"alternate_beam_head_light"`
	High_beam_head_light                int     `json:"high_beam_head_light"`
}

type Endpoint_obd2_data struct {
	Timestamp int64     `json:"timestamp"`
	Device_id string    `json"device_id"`
	Data      ODB2_data `json:"data"`
}

type MAM_send_innder struct {
	Seed    string `json:"seed"`
	Message string `json:"message"`
}

type MAM_post_send struct {
	X_API_KEY string          `json:"x-api-key"`
	Data      MAM_send_innder `json:"data"`
	Protocol  string          `json:"protocol"`
}

type MAM_send_response struct {
	Bundle_hash string `json:"bundle_hash"`
	Chid        string `json:"chid"`
	Msg_id      string `json:"msg_id"`
}

func MAM_send(host string, seed string, message string) *MAM_send_response {

	var data = MAM_post_send{
		Data: MAM_send_innder{
			Seed:    seed,
			Message: message,
		},
		Protocol: "MAM_V1",
	}

	b, _ := json.Marshal(data)

	resp, err := http.Post(host+"/mam/send", "application/json", bytes.NewBuffer(b))

	if err != nil {
		log.Println(err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Non-OK HTTP status: ", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Error message: ", string(body))
		return nil
	}

	var mam_response MAM_send_response
	json.NewDecoder(resp.Body).Decode(&mam_response)

	log.Printf("Bundle_hash:%s\n", mam_response.Bundle_hash)
	log.Printf("Chid:%s\n", mam_response.Chid)
	log.Printf("Msg id:%s\n", mam_response.Msg_id)

	return &mam_response
}

func Generate_obd2_data() ODB2_data {
	return ODB2_data{
		Vin:                                 "BCKHQEUPEEQNMKGUJ",
		Engine_load:                         rand.Intn(100),
		Engine_coolant_temperature:          rand.Intn(100),
		Fuel_pressure:                       rand.Intn(100),
		Engine_speed:                        float32(rand.Intn(100)) + rand.Float32(),
		Vehicle_speed:                       rand.Intn(100),
		Intake_air_temperature:              rand.Intn(100),
		Mass_air_flow:                       rand.Intn(100),
		Fuel_tank_level_input:               rand.Intn(100),
		Absolute_barometric_pressure:        rand.Intn(100),
		Control_module_voltage:              float32(rand.Intn(100)) + rand.Float32(),
		Throttle_position:                   rand.Intn(100),
		Ambient_air_temperature:             rand.Intn(100),
		Relative_accelerator_pedal_position: rand.Intn(100),
		Engine_oil_temperature:              rand.Intn(100),
		Engine_fuel_rate:                    float32(rand.Intn(100)) + rand.Float32(),
		Service_distance:                    rand.Intn(100),
		Anti_lock_barking_active:            rand.Intn(100),
		Steering_wheel_angle:                rand.Intn(100),
		Position_of_doors:                   rand.Intn(100),
		Right_left_turn_signal_light:        rand.Intn(100),
		Alternate_beam_head_light:           rand.Intn(100),
		High_beam_head_light:                rand.Intn(100),
	}

}

func Aes256(plaintext string, key string, iv string, blockSize int) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Padding([]byte(plaintext), blockSize, len(plaintext))
	block, _ := aes.NewCipher(bKey)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext)
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

type Endpoint_serial struct {
	IV             string
	Timestamp      int64
	Hmac           string
	Ciphertext_len int
	Ciphertext     string
}

func Endpoint_serializer(iv string, timestamp int64, hmac string, message string) string {
	endpoint_serial := fmt.Sprintf("%016s%020d%032s%010d%s", iv, timestamp, hmac, len(message), message)
	return endpoint_serial
}

func main() {
	var seed string = "BLPTDKBWYFLZRONYSSDDDEQYOHNTFUPWQVQFQGRCFPOQMIIDDVBZXQBSKTGAPESIZQKZRPOVVVOABHTTT"
	var private_key string = "LLHRCBHHYWKAGXMYCEKJIPBATQZPBQIE"
	var hmac = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	var iv string = "1234567890123456"
	var device_id string = "PWFJOIOZGKUIOBVY"
	var timestamp = time.Now().Unix()

	for i := 1; i <= 10; i++ {
		obd2_data := Generate_obd2_data()

		data := &Endpoint_obd2_data{
			Timestamp: time.Now().Unix(),
			Device_id: device_id,
			Data:      obd2_data,
		}

		data_json, _ := json.Marshal(data)

		fmt.Println(string(data_json))

		ciphertext := Aes256(string(data_json), private_key, iv, aes.BlockSize)
		var msg = Endpoint_serializer(string(iv), timestamp, hmac, string(ciphertext))
		fmt.Println(msg)
		MAM_send("http://node1.puyuma.org:6666", seed, msg)
		fmt.Println("Finish :", i)
	}
}
