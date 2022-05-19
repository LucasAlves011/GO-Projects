/**
@author : Lucas Alves - Disciplina de Paradigmas de Programação
Minha ideia foi a seguinte, 10 sensores mandam sinais para uma central usando um número aleatório como critério (mais explicações quanto a isso está abaixo).
Se um sinal for positivo, o alarme tocará por 5 segundos(independente dos sinais dos outros sensores) e voltará a receber os sinais dos sensores.
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Yellow        = "\u001b[33m"
	Reset         = "\u001b[0m"
	ChanceDeTocar = 4 // A chance em porcentagem de um sensor detectar um movimento
)

//Tipo central e seu "construtor"
type Central struct {
	ativado bool
	sinal   chan bool
	tocando bool
}

func newCentral(ativado bool, canal chan bool) Central {
	alarme := Central{ativado, canal, false}
	go alarme.esperarSinal() // Assim que uma central é criada, uma thread é lançada e a central receberá os sinais dos sensores indefinidamente
	return alarme
}

//Tipo Sensor e seu "construtor"
type Sensor struct {
	name string
}

func newSensor(name string, central Central) Sensor {
	sensor := Sensor{name}
	go sensor.mandarSinal(central) // Assim que um sensor é criado, uma thread é iniciado e começa a mandar sinais para central
	return sensor
}

//Manda um sinal para a central
func (s Sensor) mandarSinal(central Central) {
	for {
		if rand.Intn(100) < ChanceDeTocar { // Um número aleatório de 0 a 100 é gerado e se for menor que a ChanceDeTocar, manda um sinal positivo pelo canal, caso contrário manda um sinal negativo
			central.sinal <- true
		} else {
			central.sinal <- false
		}
		time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond) // Dorme por um período de 0 a 1,5 segundo
	}
}

// Toca um sinal por 5 segundos
func (c Central) tocarAlarmePor5Segundos(tocando *bool) {
	for i := 0; i < 5; i++ {
		fmt.Printf(string(Yellow)+"%d segundos - Tocando alarme...   "+Reset, i+1)
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	*tocando = false
}

func (a Central) esperarSinal() {
	for {
		if !a.tocando {
			if <-a.sinal && a.ativado {
				a.tocando = true
				a.tocarAlarmePor5Segundos(&a.tocando)
			} else {
				fmt.Println("Nada acontecendo") // Lembrando que essa impressão do nada acontecendo acontece em clocks de 50 mili segundos, e a do alarme acontece em clocks de 1 segundo. Resumindo 20 mensagens do "nada acontecendo" equivalem ao tempo de uma mensagem do "alarme".
			}
			time.Sleep(50 * time.Millisecond) // A central não pode ignorar sinais vindo dos sensores por isso ela dorme por tão pouco tempo, o ideal seria ela não dormir, mas para melhorar o efeito da visualização no console coloquei um pequeno tempo.
		}
	}
}

func main() {
	canal := make(chan bool)
	a1 := newCentral(true, canal)

	newSensor("Garagem", a1)
	newSensor("Cozinha", a1)
	newSensor("Quarto 1", a1)
	newSensor("Quarto 2", a1)
	newSensor("Dispensa", a1)
	newSensor("Garagem2", a1)
	newSensor("Dispensa", a1)
	newSensor("Banheiro???!", a1)
	newSensor("Escritório", a1)
	newSensor("Área de Serviço", a1)

	time.Sleep(time.Second * 120) // Essa simulação ficará em execução por 2 minutos, mas pode ser alterado sem problemas

}
