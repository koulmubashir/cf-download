package main_test

import (
	. "github.com/ibmjstart/cf-download"
	"github.com/ibmjstart/cf-download/cmd_exec/cmd_exec_fake"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// unit tests of individual functions
var _ = Describe("CfDownload", func() {
	Describe("Test ParseArgs functionality", func() {

		Context("Check if overWrite flag works", func() {
			It("Should set the overwrite_flag", func() {
				args := [...]string{"download", "app", "app/files/htdocs", "--overwrite"}

				flagVals, _ := ParseArgs(args[:])
				Expect(flagVals.OverWrite_flag).To(BeTrue())
				Expect(flagVals.File_flag).To(BeFalse())
				Expect(flagVals.Instance_flag).To(Equal("0"))
				Expect(flagVals.Verbose_flag).To(BeFalse())
				Expect(flagVals.Omit_flag).To(Equal(""))
			})
		})

		Context("Check if file flag works", func() {
			It("Should set the file_flag", func() {
				args := [...]string{"download", "app", "--file"}

				flagVals, _ := ParseArgs(args[:])
				Expect(flagVals.OverWrite_flag).To(BeFalse())
				Expect(flagVals.File_flag).To(BeTrue())
				Expect(flagVals.Instance_flag).To(Equal("0"))
				Expect(flagVals.Verbose_flag).To(BeFalse())
				Expect(flagVals.Omit_flag).To(Equal(""))
			})
		})

		Context("Check if verbose flag works", func() {
			It("Should set the verbose_flag", func() {
				args := [...]string{"download", "app", "--verbose"}

				flagVals, _ := ParseArgs(args[:])
				Expect(flagVals.OverWrite_flag).To(BeFalse())
				Expect(flagVals.File_flag).To(BeFalse())
				Expect(flagVals.Instance_flag).To(Equal("0"))
				Expect(flagVals.Verbose_flag).To(BeTrue())
				Expect(flagVals.Omit_flag).To(Equal(""))
			})
		})

		Context("Check if instance (i) flag works", func() {
			It("Should set the instance_flag", func() {
				args := [...]string{"download", "app", "--i", "3"}

				flagVals, _ := ParseArgs(args[:])
				Expect(flagVals.OverWrite_flag).To(BeFalse())
				Expect(flagVals.File_flag).To(BeFalse())
				Expect(flagVals.Instance_flag).To(Equal("3"))
				Expect(flagVals.Verbose_flag).To(BeFalse())
				Expect(flagVals.Omit_flag).To(Equal(""))
			})
		})

		Context("Check if omit flag works", func() {
			It("Should set the omit_flag", func() {
				args := [...]string{"download", "app", "--omit", "app/node_modules"}

				flagVals, _ := ParseArgs(args[:])
				Expect(flagVals.OverWrite_flag).To(BeFalse())
				Expect(flagVals.File_flag).To(BeFalse())
				Expect(flagVals.Instance_flag).To(Equal("0"))
				Expect(flagVals.Verbose_flag).To(BeFalse())
				Expect(flagVals.Omit_flag).To(Equal("app/node_modules"))
			})
		})

		Context("Check if correct number of paths are returned", func() {
			It("Should return 0 paths", func() {
				args := [...]string{"download", "app"}

				_, paths := ParseArgs(args[:])
				Expect(len(paths)).To(Equal(0))
			})

			It("Should return 1 path", func() {
				args := [...]string{"download", "app", "path/to/file"}

				_, paths := ParseArgs(args[:])
				Expect(len(paths)).To(Equal(1))
			})

			It("Should return 2 paths", func() {
				args := [...]string{"download", "app", "path/to/file", "path/to/other/file"}

				_, paths := ParseArgs(args[:])
				Expect(len(paths)).To(Equal(2))
			})
		})
	})

	Describe("test directoryContext parsing", func() {

		It("Should return correct strings", func() {
			paths := [...]string{"app/src/node"}

			currentDirectory, _ := os.Getwd()
			currentDirectory = filepath.ToSlash(currentDirectory)
			pathVals := GetDirectoryContext(currentDirectory, paths[:], false)

			correctSuffix := strings.HasSuffix(pathVals[0].RootWorkingDirectoryLocal, filepath.FromSlash("/cf-download/node/"))
			Expect(correctSuffix).To(BeTrue())

			Expect(pathVals[0].StartingPathServer).To(Equal("/app/src/node/"))
		})

		It("should still return /app/src/node/ for startingPath (INPUT has leading and trailing slash)", func() {
			paths := [...]string{"/app/src/node/"}

			currentDirectory, _ := os.Getwd()
			currentDirectory = filepath.ToSlash(currentDirectory)
			pathVals := GetDirectoryContext(currentDirectory, paths[:], false)

			correctSuffix := strings.HasSuffix(pathVals[0].RootWorkingDirectoryLocal, filepath.FromSlash("/cf-download/node/"))
			Expect(correctSuffix).To(BeTrue())

			Expect(pathVals[0].StartingPathServer).To(Equal("/app/src/node/"))
		})

		It("should still return /app/src/node/ for startingPath (INPUT only has trailing slash)", func() {
			paths := [...]string{"app/src/node/"}

			currentDirectory, _ := os.Getwd()
			currentDirectory = filepath.ToSlash(currentDirectory)
			pathVals := GetDirectoryContext(currentDirectory, paths[:], false)

			correctSuffix := strings.HasSuffix(pathVals[0].RootWorkingDirectoryLocal, filepath.FromSlash("/cf-download/node/"))
			Expect(correctSuffix).To(BeTrue())

			Expect(pathVals[0].StartingPathServer).To(Equal("/app/src/node/"))
		})

		It("should still return /app/src/node/ for startingPath (INPUT only has leading slash)", func() {
			paths := [...]string{"/app/src/node"}

			currentDirectory, _ := os.Getwd()
			currentDirectory = filepath.ToSlash(currentDirectory)
			pathVals := GetDirectoryContext(currentDirectory, paths[:], false)

			correctSuffix := strings.HasSuffix(pathVals[0].RootWorkingDirectoryLocal, filepath.FromSlash("/cf-download/node/"))
			Expect(correctSuffix).To(BeTrue())

			Expect(pathVals[0].StartingPathServer).To(Equal("/app/src/node/"))
		})

		It("should return /app/src/file.html for startingPath (--file flag specified)", func() {
			paths := [...]string{"/app/src/file.html"}

			currentDirectory, _ := os.Getwd()
			currentDirectory = filepath.ToSlash(currentDirectory)
			pathVals := GetDirectoryContext(currentDirectory, paths[:], true)

			correctSuffix := strings.HasSuffix(pathVals[0].RootWorkingDirectoryLocal, filepath.FromSlash("/cf-download/file.html"))
			Expect(correctSuffix).To(BeTrue())

			Expect(pathVals[0].StartingPathServer).To(Equal("/app/src/file.html"))
		})

		It("should return two staringPaths, /app/src/ and /app/logs/", func() {
			paths := [...]string{"/app/src/", "app/logs/"}

			currentDirectory, _ := os.Getwd()
			currentDirectory = filepath.ToSlash(currentDirectory)
			pathVals := GetDirectoryContext(currentDirectory, paths[:], false)

			correctSuffix := strings.HasSuffix(pathVals[0].RootWorkingDirectoryLocal, filepath.FromSlash("/cf-download/src/"))
			Expect(correctSuffix).To(BeTrue())
			correctSuffix = strings.HasSuffix(pathVals[1].RootWorkingDirectoryLocal, filepath.FromSlash("/cf-download/logs/"))
			Expect(correctSuffix).To(BeTrue())

			Expect(pathVals[0].StartingPathServer).To(Equal("/app/src/"))
			Expect(pathVals[1].StartingPathServer).To(Equal("/app/logs/"))
		})

	})

	Describe("test expandGlobs parsing", func() {
		It("should return x.txt, y.txt, a.go and ab.go", func() {
			cmdExec := cmd_exec_fake.NewCmdExec()
			cmdExec.SetOutput("Getting files for app Test in org test / space dev as user...\nOK\n\nxyz.txt                                   220B\na.go                                      675B\nab.go                                     333B\nyz.go                                     123B\n\n")
			cmdExec.SetFakeDir(false)

			paths := make([]string, 1)
			paths[0] = "*.txt"
			paths = ExpandGlobs(cmdExec, paths, "0")
			Expect(paths[0]).To(Equal("./xyz.txt"))

			paths[0] = "?.go"
			paths = ExpandGlobs(cmdExec, paths, "0")
			Expect(paths[0]).To(Equal("./a.go"))

			paths[0] = "[a-z]b.go"
			paths = ExpandGlobs(cmdExec, paths, "0")
			Expect(paths[0]).To(Equal("./ab.go"))
		})
	})

	Describe("test error catching in run() [MUST HAVE PLUGIN INSTALLED TO PASS]", func() {
		Context("when appname begins with -- or -", func() {
			It("Should print error, because user has flags before appname", func() {
				cmd := exec.Command("cf", "download", "--appname")
				output, _ := cmd.CombinedOutput()
				Expect(strings.Contains(string(output), "Error: App name begins with '-' or '--'. correct flag usage: 'cf download APP_NAME [--flags]'")).To(BeTrue())
			})

			It("Should print error, because user not specified an appName", func() {
				cmd := exec.Command("cf", "download")
				output, _ := cmd.CombinedOutput()
				Expect(strings.Contains(string(output), "Error: Missing App Name")).To(BeTrue())
			})

			It("Should print error, test overwrite flag functionality", func() {
				// create directory that needs to be overwritten
				os.Mkdir("test", 755)

				cmd := exec.Command("cf", "download", "test")
				output, _ := cmd.CombinedOutput()

				// clean up
				os.RemoveAll("test")

				Expect(strings.Contains(string(output), "already exists.")).To(BeTrue())
			})

			It("Should print error, instance flag not int", func() {
				cmd := exec.Command("cf", "download", "test", "-i", "hello")
				output, _ := cmd.CombinedOutput()
				Expect(strings.Contains(string(output), "Error:  invalid value ")).To(BeTrue())
			})

			It("Should print error, invalid flag", func() {
				cmd := exec.Command("cf", "download", "test", "-ooverwrite")
				output, _ := cmd.CombinedOutput()
				Expect(strings.Contains(string(output), "Error:  flag provided but not defined: -ooverwrite")).To(BeTrue())
			})
		})
	})

})
