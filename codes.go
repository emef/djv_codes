package djv_codes

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"sync"
)

type CodeManager struct {
	codes []string
	cursor int
  mutex *sync.Mutex
	usedCodes *os.File
}

func NewCodeManager(unusedCodesDir string, usedCodesPath string) (*CodeManager, error) {
	usedCodesMap := make(map[string]bool)
	if usedCodes, err := readCodes(usedCodesPath); err == nil {
		for _, code := range usedCodes {
			usedCodesMap[code] = true
		}
	}

	files, err := ioutil.ReadDir(unusedCodesDir)
	if err != nil {
		return nil, err
	}

	availCodesMap := make(map[string]bool, 0)
	for _, file := range files {
		if !file.IsDir() {
			path := path.Join(unusedCodesDir, file.Name())
			fileCodes, err := readCodes(path)
			if err != nil {
				return nil, err
			}
			for _, code := range fileCodes {
				_, found := usedCodesMap[code]
				if !found {
					availCodesMap[code] = true
				}
			}
		}
	}

	codes := make([]string, 0, len(availCodesMap))
	for code := range availCodesMap {
		codes = append(codes, code)
	}

	sort.Strings(codes)

	usedCodesFile, err := os.OpenFile(usedCodesPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}

	manager := &CodeManager{
		codes: codes,
		cursor: 0,
		mutex: &sync.Mutex{},
		usedCodes: usedCodesFile}

	return manager, nil
}

func readCodes(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	codes := make([]string, 0)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			codes = append(codes, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	} else {
		return codes, err
	}
}

func (manager *CodeManager) NextCode() (string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if len(manager.codes) == 0 || manager.cursor >= len(manager.codes) {
		return "", errors.New("No codes available")
	}

	code := manager.codes[manager.cursor]
	_, err := manager.usedCodes.WriteString(code + "\n")
	manager.usedCodes.Sync()

	if err != nil {
		return "", err
	} else {
		manager.cursor += 1
		return code, nil
	}
}
