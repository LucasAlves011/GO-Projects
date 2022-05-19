/**
@author : Lucas Alves - Disciplina de Paradigmas de Programação
Minha ideia foi a seguinte, ter uma variável global que controla o fluxo das threads. Já que essa variável somente é usada
para leitura por todas as threads.
A função gui() é responsável pela interação com o usuário. É nela onde a ocorre a mudança na variável estado
*/

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/jackdanger/collectlinks"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//termo a ser buscado
var term string

const (
	Amarelo  = "\u001b[33m"
	Reset1   = "\u001b[0m"
	Vermelho = "\u001b[31m"
	Verde    = "\u001b[32m"
)

var sliceString = make([]string, 0, 2) // Um slice (uma espécie de lista da linguagem) para guardar as informações das páginas (OBS: Eu tava no meio da implementação de lista encadeada com ponteiros até descobrir a existência disso aqui, mt obrigado GO <3 xD)

//mapa de paginas ja visitadas
var estado = 2 // variável que controla o fluxo das threads. 2 é o estado de execução
var visited = make(map[string]bool)

func main() {
	flag.Parse()

	//verificacao dos argumentos do main (URL inicial e termo a ser buscado)
	args := flag.Args()
	fmt.Println(args)
	if len(args) < 2 {
		fmt.Println("Por favor especifique a página inicial e o termo de pesquisa")
		os.Exit(1)
	} else if args[1] == "" {
		fmt.Println("Por favor especifique o termo de pesquisa")
		os.Exit(1)
	}

	term = args[1]

	go gui()
	// fila de URLs
	queue := make(chan string)

	go func() { queue <- args[0] }()

	for uri := range queue {
		if !visited[uri] {
			enqueue(uri, queue)
		}
	}

	fmt.Println("terminou")
}

//busca pelo termo na pagina e busca links para continuar a navegacao
func enqueue(uri string, queue chan string) {
	//fmt.Println("fetching", uri)
	visited[uri] = true

	for estado == 1 { // Enquanto estado == 1 a thread dormirá por 2 segundos, até que a variável estado deixe de ser igual a 1
		time.Sleep(2 * time.Second)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{Transport: transport}

	if !visited[uri] {
		enqueue(uri, queue)
	}

	go searchTerm(uri, client)

	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	//busca os links na pagina atual e enfileira atraves do envio no canal 'queue'
	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := fixUrl(link, uri)
		if uri != "" {
			if !visited[absolute] {
				go func() { queue <- absolute }()
			}
		}
	}
}

//funcao para procurar se o termo existe na pagina 'uri'
func searchTerm(uri string, client http.Client) {

	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	html, errr := ioutil.ReadAll(resp.Body)
	if errr != nil {
		panic(errr)
	}
	s := string(html[:len(html)])

	numero := strings.Count(strings.ToLower(s), strings.ToLower(term))                                          // Nessa linha é feita a contagem de quantas palavras existem na página
	sliceString = append(sliceString, fmt.Sprintf("Em %s a palavra %s apareceu %d vezes\n", uri, term, numero)) // Nessa linha é adicionado a lista uma string com as informações da página (uri e número de vezes que o termo aparece)
}

//conserta URLs com problemas
func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

func gui() { // Aqui é o método onde é possível mudar o estado do algoritmo
	for {
		var escolha = 0

		if estado == 1 { //Quando o crawler estiver pausado
			fmt.Println("1 - Pausar crawler" + Vermelho + " (Pausado)" + Reset)
			fmt.Println("2 - Continuar crawler")
		} else { // Em execução
			fmt.Println("1 - Pausar crawler " + Verde + " (Em execução)" + Reset)
		}

		fmt.Println("3 - Finalizar")
		fmt.Println("4 - Listar o ranking de páginas")
		fmt.Scan(&escolha)

		if escolha == 1 { // Pausa o crawler
			estado = 1
			fmt.Println(Amarelo + "Crawler pausado..." + Reset1)
		} else if escolha == 2 { // Executa o crawler
			estado = 2
			fmt.Println("Crawler reiniciado...")
		} else if escolha == 3 { // Finaliza o programa
			os.Exit(1)
		} else if escolha == 4 { // Mostra o rank
			fmt.Println(sliceString)
			fmt.Print("Enter para voltar ao menu...")
			fmt.Scanln()
		}
	}

}
