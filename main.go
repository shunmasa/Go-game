package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// GameState represents different game states
type GameState int

const (
	Start GameState = iota
	Room1
	Room2
	Room3
	Win
	GameOver
	Purchase
)

// Item represents different items in the game
type Item int

const (
	Key  Item = iota
	Coin Item = iota
)

func main() {
	inventory, err := loadGame()
	if err != nil {
		fmt.Println("Error loading game:", err)
		return
	}

	playGame(Start, &inventory)
}

func playGame(gameState GameState, inventory *[]Item) {
	switch gameState {
	case Start:
		for {
			coinCount := getCoinCount(inventory)
			fmt.Println("Welcome to the Text Adventure Game!")
			fmt.Println("You find yourself in a dark room. There are three doors in front of you.")
			fmt.Printf("You have %d coins.\n", coinCount)
			fmt.Println("Choose a door to enter (1, 2, 3):")

			var input string
			fmt.Scanln(&input)

			switch input {
			case "1", "2", "3":
				playGame(randomGameState(), inventory)
			default:
				fmt.Println("Invalid choice! You stumble in the darkness.")
				playGame(GameOver, inventory)
			}
		}

	case Room1:
		fmt.Println(getRandomRoom1Message())
		if hasKey(inventory) {
			fmt.Println("You use the key to unlock the door. The door creaks open.")
			playGame(Win, inventory)
		} else {
			fmt.Println("You search the room, trying to find the key.")
			playGame(randomGameState(), inventory)
		}

	case Room2:
		fmt.Println("Room 2 reveals itself to you. A mysterious table is adorned with a key and a coin.")
	
		var input string
		fmt.Print("What will you do? (Type 'pick up' or 'leave'): ")
		fmt.Scanln(&input)
	
		switch input {
		case "pick up":
			fmt.Println("You picked up the key and the coin. The room shivers.")
			*inventory = append(*inventory, Key, Coin)
			playGame(randomGameState(), inventory)
		case "leave":
			fmt.Println("You decide to leave the key and the coin on the table. The room remains still.")
			playGame(randomGameState(), inventory)
		default:
			fmt.Println("Invalid choice! The room reacts strangely.")
			playGame(GameOver, inventory)
		}

	case Room3:
		fmt.Println(getRandomRoom3Message())
		playGame(randomGameState(), inventory)

	case Win:
		fmt.Println("Congratulations! You unlocked the door and won the game.")

	case GameOver:
		fmt.Println("Game Over! You made a wrong choice. The darkness consumes you.")

	case Purchase:
		coinCount := getCoinCount(inventory)

		fmt.Println("Welcome to the item shop!")
		if coinCount > 0 {
			fmt.Printf("You have %d coins.\n", coinCount)
		} else {
			fmt.Println("Player has 0 coins.")
		}

		if coinCount >= 5 {
			fmt.Println("You can purchase a key to win the game (K - 5 coins).")
		}

		fmt.Println("Choose an item to purchase (1. Coin - 2 coins, 2. Back):")

		for {
			var input string
			fmt.Scanln(&input)

			switch input {
			case "1":
				if coinCount >= 2 {
					fmt.Println("You purchased a coin! The shopkeeper nods.")
					*inventory = append(*inventory, Coin)
					// Consume 2 coins
					*inventory = filterCoins(*inventory)
					*inventory = append(*inventory, Coin)
				} else {
					fmt.Println("Not enough coins to purchase the coin. The shopkeeper frowns.")
				}
			case "K":
				if coinCount >= 5 {
					fmt.Println("You purchased a key and won the game! The universe bends to your will.")
					playGame(Win, inventory)
					return
				} else {
					fmt.Println("Not enough coins to purchase the key. The shopkeeper shakes his head.")
				}
			case "2":
				playGame(randomGameState(), inventory)
			default:
				fmt.Println("Invalid choice! The shopkeeper looks confused.")
			}
		}

	default:
		fmt.Println("Do you want to save the game? (yes/no):")

		for {
			var input string
			fmt.Scanln(&input)

			switch input {
			case "yes":
				keyCode := generateKeyCode()
				err := saveGameWithKey(inventory, keyCode)
				if err != nil {
					fmt.Println("Failed to save the game:", err)
				} else {
					fmt.Printf("Game saved with key code: %s. The universe remembers.\n", keyCode)
				}
			case "no":
				fmt.Println("Thanks for playing! The adventure ends here.")
			default:
				fmt.Println("Invalid choice! The universe is indifferent.")
			}
		}
	}
}

func randomGameState() GameState {
	rand.Seed(time.Now().UnixNano())
	return GameState(rand.Intn(3))
}

func generateKeyCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func saveGameWithKey(inventory *[]Item, keyCode string) error {
	file, err := os.Create("save_game.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Key Code: %s\n", keyCode)
	if err != nil {
		return err
	}

	for _, item := range *inventory {
		switch item {
		case Key:
			_, err = fmt.Fprintln(file, "Key")
		case Coin:
			_, err = fmt.Fprintln(file, "Coin")
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func loadGame() ([]Item, error) {
	content, err := ioutil.ReadFile("save_game.txt")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 {
		return nil, nil
	}

	var inventory []Item
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "Key" {
			inventory = append(inventory, Key)
		} else if line == "Coin" {
			inventory = append(inventory, Coin)
		}
	}

	return inventory, nil
}

func getCoinCount(inventory *[]Item) uint {
	var count uint
	for _, item := range *inventory {
		if item == Coin {
			count++
		}
	}
	return count
}

func hasKey(inventory *[]Item) bool {
	for _, item := range *inventory {
		if item == Key {
			return true
		}
	}
	return false
}

func filterCoins(inventory []Item) []Item {
	var result []Item
	for _, item := range inventory {
		if item != Coin {
			result = append(result, item)
		}
	}
	return result
}

func getRandomRoom1Message() string {
	room1Messages := []string{
		"You enter Room 1. It's dark and musty. A mysterious sound echoes. You need to find a key to unlock the door.",
		"You step into Room 1. The air is heavy, and you can feel a presence. Find the key to proceed.",
		"Room 1 welcomes you with darkness. Your only way out is to uncover the key hidden within.",
	}

	return room1Messages[rand.Intn(len(room1Messages))]
}

// func getRandomRoom2Message() string {
// 	room2Messages := []string{
// 		"You enter Room 2. It's dimly lit with a strange aura. A table stands in the center.",
// 		"On the table, there's a key and a coin. Do you want to pick them up? (yes/no):",
// 		"Room 2 reveals itself to you. A mysterious table is adorned with a key and a coin. What will you do?",
// 	}

// 	return room2Messages[rand.Intn(len(room2Messages))]
// }

func getRandomRoom3Message() string {
	room3Messages := []string{
		"You enter Room 3. A giant spider blocks your way! You can't proceed this way. Go back to another room.",
		"A massive spider guards Room 3. Retreat to another room to escape its web.",
		"Room 3 presents a challenge - a giant spider. Your only option is to turn back and explore another path.",
	}

	return room3Messages[rand.Intn(len(room3Messages))]
}
