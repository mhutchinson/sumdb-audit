package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/mod/sumdb/note"
	"golang.org/x/mod/sumdb/tlog"
)

var (
	height = flag.Int("h", 8, "tile height")
	offset = flag.Int("o", 0, "offset")
	vkey   = flag.String("k", "sum.golang.org+033de0ae+Ac4zctda0e5eza+HJyk9SxEdh+s3Ux18htTTAD8OuAn8", "key")
	url    = flag.String("u", "", "url to server (overriding name)")
	vflag  = flag.Bool("v", false, "enable verbose output")
)

func main() {
	log.SetPrefix("audit: ")
	log.SetFlags(0)
	flag.Parse()

	size, root := getLatestCheckpoint()
	log.Printf("tree size = %d, root hash = %s", size, root)
	path := dataPathForOffset(*height, *offset)
	data, err := readRemote(path)
	if err != nil {
		log.Fatalf("failed to read %s: %s", path, err)
	}
	log.Printf("%s", data)
}

func getLatestCheckpoint() (int64, tlog.Hash) {
	checkpoint, err := readRemote("/latest")
	if err != nil {
		log.Fatalf("failed to get /latest Checkpoint: %s", err)
	}

	verifier, err := note.NewVerifier(*vkey)
	if err != nil {
		log.Fatal(err)
	}
	verifiers := note.VerifierList(verifier)

	note, err := note.Open(checkpoint, verifiers)
	if err != nil {
		log.Fatal(err)
	}
	tree, err := tlog.ParseTree([]byte(note.Text))
	if err != nil {
		log.Fatal(err)
	}

	return tree.N, tree.Hash
}

func readRemote(path string) ([]byte, error) {
	name := *vkey
	if i := strings.Index(name, "+"); i >= 0 {
		name = name[:i]
	}
	start := time.Now()
	target := "https://" + name + path
	if *url != "" {
		target = *url + path
	}
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET %v: %v", target, resp.Status)
	}
	data, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if *vflag {
		fmt.Fprintf(os.Stderr, "%.3fs %s\n", time.Since(start).Seconds(), target)
	}
	return data, nil
}

const pathBase = 1000

func dataPathForOffset(height, offset int) string {
	nStr := fmt.Sprintf("%03d", offset%pathBase)
	for offset >= pathBase {
		offset /= pathBase
		nStr = fmt.Sprintf("x%03d/%s", offset%pathBase, nStr)
	}
	return fmt.Sprintf("/tile/%d/data/%s", height, nStr)
}
