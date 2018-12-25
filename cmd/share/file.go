// SPDX-License-Identifier: Apache-2.0

package share

import (
	"bufio"
	"os"
	"fmt"
)

func PathExists(path string) bool {
	_, r := os.Stat(path)
	if os.IsNotExist(r) {
		return false
	}

	return true
}

func ReadFullFile(path string) (lines []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()

	return lines, nil
}

func WriteFullFile(path string, lines[] string) (error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	w.Flush()

	return nil
}

func ReadOneLineFile(path string) (string, error) {
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	line := scanner.Text()

	err = scanner.Err()

	return line, nil
}

func WriteOneLineFile(path string, line string) (error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, line)

	return w.Flush()
}

func CreateDirectory(directoryPath string, perm os.FileMode) (error) {
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(directoryPath, perm)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateDirectoryNested(directoryPath string, perm os.FileMode) (error) {
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(directoryPath, perm)
		if err != nil {
			return err
		}
	}

	return nil
}
