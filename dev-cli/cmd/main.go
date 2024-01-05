/*
Ferramenta CLI protótipo - os trechos serão comentados para esclarecimento
das funcionalidas
*/
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"strings"
)



func main() {

	var rootCommand = &cobra.Command{}
	var projectName, projectPath string

	var cmd = &cobra.Command{
		Use:   "newrest",
		Short: "Cria um boilerplate para um novo projeto api rest",
		Run: func(cmd *cobra.Command, args []string) {
			if projectName == "" {
				color.Red.Println("Você deve informar o nome do projeto.")
				return
			}
			if projectPath == "" {
				color.Red.Println("Voce deve informar pasta/caminho do projeto.")
				return
			}
			color.Cyan.Println("\nCriando projeto...\n")

			globalPath := filepath.Join(projectPath, projectName)

			if _, err := os.Stat(globalPath); err == nil {
				color.Yellow.Println("Pasta do projeto já existe.")
				return
			}
			if err := os.Mkdir(globalPath, os.ModePerm); err != nil {
				msgErr:=fmt.Sprintf("Não foi possíveL criar a pasta do projeto em %s.\n Erro: %s\n",globalPath,err)
        color.Red.Println("ERRO")
        log.Fatal(msgErr)
			}

			startGo := exec.Command("go", "mod", "init", projectName)
			startGo.Dir = globalPath
			startGo.Stdout = os.Stdout
			startGo.Stderr = os.Stderr
			err := startGo.Run()
			if err != nil {
				log.Fatal(err)
			}

			cmdPath := filepath.Join(globalPath, "cmd")
			if err := os.Mkdir(cmdPath, os.ModePerm); err != nil {
				log.Fatal(err)
			}
			internalPath := filepath.Join(globalPath, "internal")
			if err := os.Mkdir(internalPath, os.ModePerm); err != nil {
				log.Fatal(err)
			}
			handlerPath := filepath.Join(internalPath, "handler")
			if err := os.Mkdir(handlerPath, os.ModePerm); err != nil {
				log.Fatal(err)
			}
			routesPath := filepath.Join(handlerPath, "routes")
			fmt.Println(routesPath)
			if err := os.Mkdir(routesPath, os.ModePerm); err != nil {
				log.Fatal(err)
			}

			mainPath := filepath.Join(cmdPath, "main.go")
			mainFile, err := os.Create(mainPath)
			if err != nil {
				log.Fatal(err)
			}
			defer mainFile.Close()
			if err := WriteMainFile(mainPath, projectName); err != nil {
				log.Fatal(err)
			}

			routesFilePath := filepath.Join(routesPath, "routes.go")
			routesFile, err := os.Create(routesFilePath)
			if err != nil {
				log.Fatal(err)
			}
			defer routesFile.Close()
			if err := WriteRoutesFile(routesFilePath); err != nil {
				log.Fatal(err)
			}

      color.Green.Println("\nProjeto criado com sucesso!!!\n")
		},
	}

	cmd.Flags().StringVarP(&projectName, "name", "n", "", "Nome do projeto")
	cmd.Flags().StringVarP(&projectPath, "path", "p", "", "Caminho da pasta que será criado o projeto")

	rootCommand.AddCommand(cmd)
	rootCommand.Execute()
}

func WriteMainFile(mainPath string, projectName string) error {
	
  importRootFolder:= fmt.Sprintf("%s/internal/handler/routes",projectName)

  packageContent := []byte(`package main
  
  import (
      "fmt"
      "log"
      "net/http"
      "%s"
  )
  var port string
  func main() {
    fmt.Println("Iniciando o servidor REST")
    HandleRequest()
  }

  func HandleRequest(){
    //CONFIGURE A PORTA DO SERVIDOR AQUI
    port=":8080"

    //DEFINIR AS ROTAS
    http.HandleFunc("/",routes.Home)

    fmt.Printf("Servidor ativo na porta token",port)
    log.Fatal(http.ListenAndServe(port,nil))
  }
  `)

	mainFile, err := os.OpenFile(mainPath, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer mainFile.Close()
	fileContent := fmt.Sprintf(string(packageContent), importRootFolder)
  fileContent=strings.Replace(fileContent, "token","%s",1)

	_, err = mainFile.Write([]byte(fileContent))
	if err != nil {
		return err
	}

	return nil
}

func WriteRoutesFile(routesFilePath string) error {
	packageContent := []byte(`package routes
import (
      "fmt"
      "net/http"
      
  )
  // SUAS ROTAS 
  func Home(w http.ResponseWriter, r *http.Request){
    
    fmt.Println(w, "Home Page")
  }
  //...
  
  `)

	routesFile, err := os.OpenFile(routesFilePath, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer routesFile.Close()

	_, err = routesFile.Write(packageContent)
	if err != nil {
		return err
	}

	return nil
}
