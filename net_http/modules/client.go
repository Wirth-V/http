//данный модуль реазлизует клиентскую часть приложения

package modules

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Client(req *flag.FlagSet, host string, port string) error {
	if req == nil {
		return fmt.Errorf("attempt to pass nil to the 'req' variable")
	}

	//Проверяет длинну и допустимость вводимых данных
	err := check(host, port)
	if err != nil {
		return err
	}

	hostPort := strings.Join([]string{host, port}, ":")

	//Определяет функционал вложенных команд
	switch req.Arg(Zero) {
	case "list":
		list := flag.NewFlagSet("list", flag.ExitOnError)

		list.Parse(req.Args()[One:])

		InfoLog.Println("GET request:")
		err := getItems("", hostPort)
		if err != nil {
			return err
		}

	/*
		case "get_from_db":
			list := flag.NewFlagSet("list", flag.ExitOnError)
			nameList := list.String("id", "", "ID")

			list.Parse(req.Args()[One:])

			//Определяет поведение кода, когда вместо флага команды `-id {id}`
			//вводится аргумент команды `{id}` без флага
			if *nameList == "" {

				if list.Args() != nil {

					//ничего не ввили
					if list.Arg(Zero) == "" {
						return fmt.Errorf("you flag is not correct: the flag was not entered")

					}

					//ввели что-то лишнее
					if list.Arg(One) != "" {
						return fmt.Errorf("you flag is not correct: an extra flag was introduced")

					}

					*nameList = list.Arg(Zero)
				} else {
					return fmt.Errorf("you flag is not correct")
				}

			}

			err := check(*nameList)
			if err != nil {
				return err
			}


			//	InfoLog.Println("GET request:")
			//	err = getItems(*nameList, hostPort)
			//	if err != nil {
			//		return err
			//	}


			// Установка соединения с базой данных
			connFerst, err = pgx.Connect(context.Background(), connString)
			if err != nil {
				return fmt.Errorf("unable to connect to database, %v", err)

			}
			defer connFerst.Close(context.Background())


	*/

	case "get":
		list := flag.NewFlagSet("list", flag.ExitOnError)
		nameList := list.String("id", "", "ID")

		list.Parse(req.Args()[One:])

		//Определяет поведение кода, когда вместо флага команды `-id {id}`
		//вводится аргумент команды `{id}` без флага
		if *nameList == "" {

			if list.Args() != nil {

				//ничего не ввили
				if list.Arg(Zero) == "" {
					return fmt.Errorf("you flag is not correct: the flag was not entered")

				}

				//ввели что-то лишнее
				if list.Arg(One) != "" {
					return fmt.Errorf("you flag is not correct: an extra flag was introduced")

				}

				*nameList = list.Arg(Zero)
			} else {
				return fmt.Errorf("you flag is not correct")
			}

		}

		err := check(*nameList)
		if err != nil {
			return err
		}

		InfoLog.Println("GET request:")
		err = getItems(*nameList, hostPort)
		if err != nil {
			return err
		}

	case "create":
		creates := flag.NewFlagSet("create", flag.ExitOnError)
		nameCreate := creates.String("name", "", "Name")

		creates.Parse(req.Args()[One:])

		// Определяет поведение кода, когда вместо флага команды `-name {Имя}`
		// вводится аргумент команды `{Имя}` без флага
		if *nameCreate == "" {

			if creates.Args() != nil {
				//ничего не ввели
				if creates.Arg(Zero) == "" {
					return fmt.Errorf("you command is not correct: the flag was not entered")

				}

				//ввели что-то лишнее
				if creates.Arg(One) != "" {
					return fmt.Errorf("you command is not correct: an extra flag was introduced")
				}

				*nameCreate = creates.Arg(Zero)
			} else {
				return fmt.Errorf("you command is not correct")
			}
		}

		err := check(*nameCreate)
		if err != nil {
			return err
		}

		InfoLog.Println("POST request:")
		newItem := Item{Name: *nameCreate}
		err = createItem(newItem, hostPort)
		if err != nil {
			return err
		}

	case "update":
		update := flag.NewFlagSet("update", flag.ExitOnError)
		nameUpdate := update.String("name", "New Item", "Name")
		idName := update.String("id", "", "Name")

		update.Parse(req.Args()[One:])

		//Определяет поведение кода, когда вместо флагов команды `-name {Имя} -id {id}`
		//вводится флаг и аргумент команды `-name {Имя} {id}`
		//При этом игнарирование флага `-name' не допустимо
		if *idName == "" {

			if update.Args() != nil {
				//ничего не ввели
				if update.Arg(Zero) == "" {
					return fmt.Errorf("you command is not correct: the flag was not entered")
				}

				//ввели что-то лишнее
				if update.Arg(One) != "" {
					return fmt.Errorf("you command is not correct: an extra flag was introduced")
				}

				*idName = update.Arg(Zero)
			} else {
				return fmt.Errorf("you command is not correct")
			}

		}

		err := check(*nameUpdate, *idName)
		if err != nil {
			return err
		}

		newUpdat := Item{ID: *idName, Name: *nameUpdate}

		InfoLog.Println("PUT request:")
		err = updateItem(*idName, newUpdat, hostPort)
		if err != nil {
			return err
		}

	case "delete":
		delete := flag.NewFlagSet("delete", flag.ExitOnError)
		idDelete := delete.String("id", "", "ID of delete")

		delete.Parse(req.Args()[One:])

		//Определяет поведение кода, когда вместо флага команды `-id {id}`
		//вводится аргумент команды `{id}` без флага
		if *idDelete == "" {

			if delete.Args() != nil {
				//ничего не ввели
				if delete.Arg(Zero) == "" {
					return fmt.Errorf("you command is not correct: the flag was not entered")
				}
				//ввели что-то лишнее
				if delete.Arg(One) != "" {
					return fmt.Errorf("you command is not correct: an extra flag was introduced")
				}

				*idDelete = delete.Arg(Zero)
			} else {
				return fmt.Errorf("you command is not correct")
			}
		}

		err := check(*idDelete)
		if err != nil {
			return err
		}

		InfoLog.Println("DELETE request:")
		err = deleteItem(*idDelete, hostPort)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("you flag is not correct")
	}
	return nil
}

func getItems(nameList string, hostPort string) error {
	var resp *http.Response
	var err error

	// Отправка GET-запроса на сервер по указанному URL.
	if nameList == "" {
		resp, err = http.Get(fmt.Sprintf("http://%s/items/", hostPort))
	} else {
		resp, err = http.Get(fmt.Sprintf("http://%s/items/%s/", hostPort, nameList))
	}
	if err != nil {
		return fmt.Errorf("error sending the GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection problems, response code: %v", resp.StatusCode)
	}

	// Читаем и конвертируем тело ответа в байты
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err

	}

	ResponseLog.Println(string(bytesResp))

	// Обработка ответа сервера.
	err = printResponse(resp)
	if err != nil {
		return fmt.Errorf("error Processing the server response: %v", err)
	}

	return nil
}

func createItem(item Item, hostPort string) error {
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/items/", hostPort), "application/json", strings.NewReader(string(itemJSON)))
	if err != nil {
		return fmt.Errorf("error when sending a POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection problems, response code: %v", resp.StatusCode)
	}

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ResponseLog.Println(string(bytesResp))

	err = printResponse(resp)
	if err != nil {
		return fmt.Errorf("error processing the server response: %v", err)
	}

	return nil
}

func updateItem(itemID string, updatedItem Item, hostPort string) error {
	itemJSON, err := json.Marshal(updatedItem)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Создание клиента для отправки PUT-запроса. bytes.NewBuffer(itemJSON)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/items/%s/", hostPort, itemID), bytes.NewReader(itemJSON))
	if err != nil {
		return fmt.Errorf("error creating a PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending the PUT request: %v", err)
	}
	defer resp.Body.Close()

	// проверяем код ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection problems, response code: %v", resp.StatusCode)
	}

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ResponseLog.Println(string(bytesResp))

	err = printResponse(resp)
	if err != nil {
		return fmt.Errorf("error processing the server response: %v", err)
	}

	return nil
}

func deleteItem(itemID string, hostPort string) error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/items/%s/", hostPort, itemID), nil)
	if err != nil {
		return fmt.Errorf("error creating a DELETE request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending the DELETE request: %v", err)
	}

	// проверяем код ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection problems, response code: %v", resp.StatusCode)
	}

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ResponseLog.Println(string(bytesResp))

	err = printResponse(resp)
	if err != nil {
		return fmt.Errorf("error processing the server response: %v", err)
	}

	return nil
}

// printResponse выводит информацию о статусе и теле ответа сервера.
func printResponse(resp *http.Response) error {
	// Чтение тела ответа в байтовый массив.
	var bodyBytes []byte
	_, err := resp.Body.Read(bodyBytes)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading the response: %v", err)
	}

	// Вывод информации о статусе и теле ответа.
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("---------------\n\n")

	return nil

}
