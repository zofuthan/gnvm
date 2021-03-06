package command

import (

	// lib
	. "github.com/Kenshin/cprint"
	"github.com/spf13/cobra"

	// local
	"gnvm/config"
	"gnvm/nodehandle"
	"gnvm/util"
)

var (
	global bool
	remote bool
)

// defind root cmd
var gnvmCmd = &cobra.Command{
	Use:   "gnvm",
	Short: "Gnvm is Node.exe version manager for Windows by GO",
	Long: `Gnvm is Node.exe version manager for Windows by GO
           like nvm. Complete documentation is available at https://github.com/kenshin/gnvm`,
	Run: func(cmd *cobra.Command, args []string) {
		// TO DO
	},
}

// sub cmd
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gnvm.exe",
	Long: `Print the version number of gnvm.exe e.g.
gnvm version
gnvm version --remote`,
	Run: func(cmd *cobra.Command, args []string) {
		// check args
		if len(args) > 0 {
			P(WARING, "'%v' no parameter, please check your input. See '%v'.\n", "gnvm version", "gnvm help version")
		}
		nodehandle.Version(remote)
	},
}

// sub cmd
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install any node.exe version",
	Long: `Install any node.exe version e.g.
gnvm install latest
gnvm install x.xx.xx y.yy.yy
gnvm install x.xx.xx --global
gnvm install npm`,
	Run: func(cmd *cobra.Command, args []string) {
		var newArgs []string

		if len(args) == 0 {
			P(ERROR, "'%v' need parameter, please check your input. See '%v'.\n", "gnvm install", "gnvm help install")
		} else {

			if global && len(args) > 1 {
				P(WARING, "when use --global must be only one parameter, e.g. '%v'. See 'gnvm install help'.\n", "gnvm install x.xx.xx --global")
			}

			if len(args) == 1 {
				if value := util.EqualAbs("npm", args[0]); value == "npm" {
					nodehandle.InstallNpm()
					return
				}
			}

			for _, v := range args {

				v = util.EqualAbs("latest", v)
				v = util.EqualAbs("npm", v)

				// check npm
				if v == "npm" {
					P(WARING, "use format error, the correct format is '%v'. See '%v'.\n", "gnvm install npm", "gnvm help install")
					continue
				}

				// check latest
				if v == config.LATEST {
					newArgs = append(newArgs, v)
					continue
				}

				// check version format
				if ok := util.VerifyNodeVersion(v); ok != true {
					P(ERROR, "%v format error, the correct format is %v. \n", v, "x.xx.xx")
				} else {
					newArgs = append(newArgs, v)
				}
			}
			nodehandle.Install(newArgs, global)
		}
	},
}

// sub cmd
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall local node.exe version",
	Long: `Uninstall local node.exe version e.g.
gnvm uninstall 0.10.28
gnvm uninstall latest
gnvm uninstall npm
gnvm uninstall 0.10.26 0.11.2 latest
gnvm uninstall ALL`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			P(ERROR, "'%v' need parameter, please check your input. See '%v'.\n", "gnvm uninstall", "gnvm help uninstall")
			return
		} else if len(args) == 1 {

			args[0] = util.EqualAbs("npm", args[0])
			args[0] = util.EqualAbs("ALL", args[0])

			if args[0] == "npm" {
				nodehandle.UninstallNpm()
				return
			}

			if args[0] == "ALL" {

				if newArr, err := nodehandle.LS(false); err != nil {
					P(ERROR, "remove all folder Error: %v\n", err.Error())
					return
				} else {
					args = newArr
				}
			}

		}

		for _, v := range args {

			v = util.EqualAbs("ALL", v)
			v = util.EqualAbs("latest", v)
			v = util.EqualAbs("npm", v)

			if v == "npm" {
				P(WARING, "use format error, the correct format is '%v'. See '%v'.\n", "gnvm uninstall npm", "gnvm help uninstall")
				continue
			}

			if v == "ALL" {
				P(WARING, "use format error, the correct format is '%v'. See '%v'.\n", "gnvm uninstall ALL", "gnvm help uninstall")
				continue
			}

			// get true version
			v = nodehandle.TransLatestVersion(v, true)

			// check version format
			if ok := util.VerifyNodeVersion(v); ok != true {
				P(ERROR, "%v format error, the correct format is %v.\n", v, "x.xx.xx")
			} else {
				nodehandle.Uninstall(v)
			}
		}
	},
}

// sub cmd
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use any version of the local already exists",
	Long: `Use any version of the local already exists e.g.
gnvm use x.xx.xx
gnvm use latest`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {

			args[0] = util.EqualAbs("latest", args[0])

			if args[0] != "latest" && util.VerifyNodeVersion(args[0]) != true {
				P(ERROR, "use parameter support '%v' or '%v', e.g. %v, please check your input. See '%v'.\n", "latest", "x.xx.xx", "0.10.28", "gnvm help use")
				return
			}

			// set use
			if ok := nodehandle.Use(args[0]); ok == true {
				config.SetConfig(config.GLOBAL_VERSION, nodehandle.TransLatestVersion(args[0], false))
			}
		} else {
			P(ERROR, "use parameter maximum is 1, please check your input. See '%v'.\n", "gnvm help use")
		}
	},
}

// sub cmd
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update latest node.exe",
	Long: `Update latest node.exe e.g.
gnvm update latest
gnvm update latest --global`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			args[0] = util.EqualAbs("latest", args[0])
			args[0] = util.EqualAbs("gnvm", args[0])
			switch args[0] {
			case "latest":
				nodehandle.Update(global)
			case "gnvm":
				P(WARING, "%v temporarily does not support. See '%v'.\n", args[0], "gnvm help update")
			default:
				P(ERROR, "gnvm update only support <%v> parameter. See '%v'.\n", "latest", "gnvm help update")
			}
		} else {
			P(ERROR, "use parameter maximum is 1, temporary support only <%v>, please check your input. See '%v'.\n", "latest", "gnvm help update")
		}
	},
}

// sub cmd
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List show all <local> <remote> node.exe version",
	Long: `List show all <local> <remote> node.exe version e.g.
gnvm ls
gnvm ls --remote`,
	Run: func(cmd *cobra.Command, args []string) {

		// check args
		if len(args) > 0 {
			P(WARING, "gnvm ls no parameter, please check your input. See '%v'.\n", "gnvm help ls")
		}

		if remote == true {
			nodehandle.LsRemote()
		} else {
			// check local ls
			nodehandle.LS(true)
		}
	},
}

// sub cmd
var nodeVersionCmd = &cobra.Command{
	Use:   "node-version",
	Short: "Show <global> <latest> node.exe version",
	Long: `Show <global> <latest> node.exe version e.g.
gnvm node-version
gnvm node-version latest
gnvm node-version latest --remote
gnvm node-version global`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			P(WARING, "use parameter maximum is 1, temporary only support <%v>, <%v>, please check your input. See '%v'.\n", "global", "latest", "gnvm help node-version")
		} else if len(args) == 1 {

			args[0] = util.EqualAbs("global", args[0])
			args[0] = util.EqualAbs("latest", args[0])

			switch {
			case args[0] != "global" && args[0] != "latest":
				P(WARING, "gnvm node-version only support <%v>, <%v> parameter.\n", "global", "latest")
			case args[0] != "latest" && remote:
				P(WARING, "gnvm node-version only support <%v> parameter.\n", "latest --remote")
			}
		}
		nodehandle.NodeVersion(args, remote)
	},
}

// sub cmd
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Setter and getter registry",
	Long: `Setter and getter registry e.g.
gnvm config registry
gnvm config registry http://dist.u.qiniudn.com/
gnvm config registry DEFAULT
gnvm config INIT`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {

			args[0] = util.EqualAbs("registry", args[0])
			args[0] = util.EqualAbs("noderoot", args[0])
			args[0] = util.EqualAbs("INIT", args[0])

			if args[0] == "INIT" {
				config.ReSetConfig()
				return
			}

			P(DEFAULT, "gnvm config %v is %v\n", args[0], config.GetConfig(args[0]))

		} else if len(args) == 2 {
			args[0] = util.EqualAbs("registry", args[0])
			args[0] = util.EqualAbs("noderoot", args[0])
			args[1] = util.EqualAbs("DEFAULT", args[1])
			switch {
			case args[0] == "registry" && args[1] != "DEFAULT":
				if newValue := config.SetConfig(args[0], args[1]); newValue != "" {
					P(DEFAULT, "Set success, %v new value is %v\n", args[0], newValue)
				}
			case args[0] == "registry" && args[1] == "DEFAULT":
				if newValue := config.SetConfig(args[0], config.REGISTRY_VAL); newValue != "" {
					P(DEFAULT, "Reset success, %v new value is %v\n", args[0], newValue)
				}
			case args[0] == "noderoot":
				P(ERROR, "%v read-only, temporarily does not support set way. See '%v'.\n", args[0], "gnvm help config")
			default:
				P(ERROR, "config parameter only include %v, please check your input. See '%v'.\n", "registry", "gnvm help config")
			}
		} else if len(args) > 2 {
			P(ERROR, "config parameter maximum is 2, please check your input. See '%v'.\n", "gnvm help config")
		}
	},
}

func init() {

	// add sub cmd to root
	gnvmCmd.AddCommand(versionCmd)
	gnvmCmd.AddCommand(installCmd)
	gnvmCmd.AddCommand(uninstallCmd)
	gnvmCmd.AddCommand(useCmd)
	gnvmCmd.AddCommand(updateCmd)
	gnvmCmd.AddCommand(lsCmd)
	gnvmCmd.AddCommand(nodeVersionCmd)
	gnvmCmd.AddCommand(configCmd)

	// flag
	installCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "set this version global version.")
	updateCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "set this version global version.")
	lsCmd.PersistentFlags().BoolVarP(&remote, "remote", "r", false, "get remote all node.js version list.")
	nodeVersionCmd.PersistentFlags().BoolVarP(&remote, "remote", "r", false, "get remote node.js latest version.")
	versionCmd.PersistentFlags().BoolVarP(&remote, "remote", "r", false, "get remote gnvm latest version.")

	// exec
	gnvmCmd.Execute()
}
