package main

import (
	"encoding/binary"
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	var (
		samplingRate  uint64
		bitDepth      uint64
		channel       uint64
		inputFilePath string
	)

	flag.Uint64Var(&samplingRate, "s", 44100, "Sampling rate")
	flag.Uint64Var(&bitDepth, "b", 16, "Bit depth")
	flag.Uint64Var(&channel, "c", 2, "Number of channel(1 : Monaural, 2 : Stereo)")
	flag.StringVar(&inputFilePath, "i", "", "Input file path")

	flag.Parse()

	if inputFilePath == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatal(err)
	}

	fileinfo, err := inputFile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	length := uint64(fileinfo.Size())

	outputFile, err := os.Create(fileinfo.Name() + ".wav")
	if err != nil {
		log.Fatal(err)
	}

	defer inputFile.Close()
	defer outputFile.Close()

	headerChunk := make([]byte, 44)
	copy(headerChunk[0:4], []byte("RIFF"))                                     // RIFF チャンク識別子
	copy(headerChunk[4:8], uintToByteArray(length+36))                         // 以降のサイズ
	copy(headerChunk[8:12], []byte("WAVE"))                                    // RIFF の種類
	copy(headerChunk[12:16], []byte("fmt "))                                   // fmt チャンク識別子
	copy(headerChunk[16:20], uintToByteArray(16))                              // fmt チャンクのサイズ(識別子除く)
	copy(headerChunk[20:22], uintToByteArray(1))                               // フォーマットID
	copy(headerChunk[22:24], uintToByteArray(channel))                         // チャンネル数
	copy(headerChunk[24:28], uintToByteArray(samplingRate))                    // サンプリングレート
	copy(headerChunk[28:32], uintToByteArray(samplingRate*bitDepth/8*channel)) // バイト/秒
	copy(headerChunk[32:34], uintToByteArray(bitDepth/8*channel))              // ブロックサイズ
	copy(headerChunk[34:36], uintToByteArray(bitDepth))                        // ビット深度
	copy(headerChunk[36:40], []byte("data"))                                   // data チャンク識別子
	copy(headerChunk[40:44], uintToByteArray(length))                          // data チャンクのサイズ(識別子除く)
	outputFile.Write(headerChunk)
	io.Copy(outputFile, inputFile)
}
func uintToByteArray(x uint64) (buf []byte) {
	buf = make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, x)
	return
}
