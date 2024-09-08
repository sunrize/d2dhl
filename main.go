package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

var (
	srcDirs = flag.String("src_dirs", "", "Comma separated list of source directories")
	dstDirs = flag.String("dst_dirs", "", "Comma separated list of destination directories")
	dest    = flag.String("dest", "", "Main destination directory")
	output  = flag.Bool("output", false, "Whether to output collected inodes")
	dry     = flag.Bool("dry", false, "Whether to copy files")
)

func init() {
	flag.Parse()
}

// InodesCache хранит кеш данных об inode файлов
type InodesCache struct {
	Inodes map[uint64]struct {
		SourcePath string
		Relative   string
	}
}

func NewInodesCache() *InodesCache {
	return &InodesCache{
		Inodes: make(map[uint64]struct {
			SourcePath string
			Relative   string
		}),
	}
}

// Add добавляет inode в кеш
func (c *InodesCache) Add(inode uint64, path string, sourcePath string) {
	relativePath, _ := filepath.Rel(sourcePath, path)
	c.Inodes[inode] = struct {
		SourcePath string
		Relative   string
	}{sourcePath, relativePath}
}

// Contains проверяет наличие inode в кеше
func (c *InodesCache) Contains(inode uint64) bool {
	_, ok := c.Inodes[inode]
	return ok
}

// GetPath возвращает полный путь к файлу по его inode
func (c *InodesCache) GetPath(inode uint64) string {
	entry, exists := c.Inodes[inode]
	if !exists {
		return ""
	}
	return filepath.Join(entry.SourcePath, entry.Relative)
}

// SourcePath возвращает путь до источника файла по его inode
func (c *InodesCache) SourcePath(inode uint64) string {
	entry, exists := c.Inodes[inode]
	if !exists {
		return ""
	}
	return entry.SourcePath
}

// RelativePath возвращает относительный путь к файлу от источника по его inode
func (c *InodesCache) RelativePath(inode uint64) string {
	entry, exists := c.Inodes[inode]
	if !exists {
		return ""
	}
	return entry.Relative
}

// LoadInodesFromDirs загружает все inode файлов из указанных каталогов
func LoadInodesFromDirs(dirs []string, cache *InodesCache) error {
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			cache.Add(info.Sys().(*syscall.Stat_t).Ino, path, dir)
			return nil
		})
		if err != nil {
			return fmt.Errorf("error walking directory %q: %w", dir, err)
		}
	}
	return nil
}

func CheckAndLinkFiles(srcCache *InodesCache, dstCache *InodesCache, mainDest string) error {
	srcKeys := make([]uint64, len(srcCache.Inodes))
	i := 0
	for k := range srcCache.Inodes {
		srcKeys[i] = k
		i++
	}
	sort.Slice(srcKeys, func(i, j int) bool { return srcKeys[i] < srcKeys[j] })

	for _, key := range srcKeys {
		if !dstCache.Contains(key) {
			srcPath := srcCache.GetPath(key)
			srcRelativePath := srcCache.RelativePath(key)
			// srcInfo, err := os.Lstat(srcPath)
			// if err != nil {
			// 	fmt.Errorf("error getting info on source file %q: %w", srcPath, err)
			// }
			linkedPath := filepath.Join(mainDest, srcRelativePath)
			if *output {
				fmt.Println("Source Path: ", srcPath)
				fmt.Println("Linked Path: ", linkedPath)
			}
			if !*dry {
				// Создаем родительские директории, если они отсутствуют
				err := os.MkdirAll(filepath.Dir(linkedPath), os.ModePerm)
				err = os.Link(srcPath, linkedPath)
				if err != nil {
					fmt.Errorf("error linking file from %q to %q: %w", srcPath, linkedPath, err)
				}
			}
		}
	}
	return nil
}

func main() {
	srcCache := NewInodesCache()
	dstCache := NewInodesCache()
	err := LoadInodesFromDirs(strings.Split(*srcDirs, ","), srcCache)
	if err != nil {
		log.Fatalf("Error loading inodes from source directories: %v\n", err)
	}
	if *dest == "" {
		fmt.Fprintln(os.Stderr, "Missing argument -dest!")
		os.Exit(1)
	}
	destDirsList := append(strings.Split(*dstDirs, ","), *dest)
	err = LoadInodesFromDirs(destDirsList, dstCache)
	if err != nil {
		log.Fatalf("Error loading inodes from destination directories: %v\n", err)
	}

	err = CheckAndLinkFiles(srcCache, dstCache, *dest)
	if err != nil {
		fmt.Println("Error checking and linking files:", err)
		return
	}
	// if *output {
	// 	for inode := range srcCache.Inodes {
	// 		fmt.Println("Source inode:", inode, "\t", srcCache.GetPath(inode))
	// 	}
	// 	for inode := range dstCache.Inodes {
	// 		fmt.Println("Destination inode:", inode, "\t", dstCache.GetPath(inode))
	// 	}
	// }
}
