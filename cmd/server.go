package cmd

import (
	"fmt"
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/abronan/valkeyrie/store/consul"
	"github.com/spf13/cobra"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/api"
	dwarkaStore "gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/store"
	"gitlab.com/vedhabhavanam/smarthome/dwarka/pkg/strings"
)

var (
	storeBackend     string
	boldDBFilePath   string
	consulHTTPAddr   string
	bindAddress      string
	httpPort         string
	storeBasePath    string
	bucketName       string
	supportedBackend = []string{string(store.BOLTDB), string(store.CONSUL)}
)

var _ = func() error {
	storeBackend = flagHackLookup("--store-backend")
	if storeBackend == "" {
		storeBackend = string(store.BOLTDB)
	}
	return nil
}()

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:           "server",
	Short:         "Start REST API server",
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !strings.Contains(supportedBackend, storeBackend) {
			return fmt.Errorf("unsupported store backend '%s'", storeBackend)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := dwarkaStore.NewStore(storeBasePath, storeBackend, bucketName, addrs()...)
		if err != nil {
			return err
		}
		server := api.NewServer(bindAddress, httpPort, store)

		err = store.RefreshUptime()
		if err != nil {
			return err
		}
		return server.ListenAndServe()
	},
}

func addrs() []string {
	backend := store.Backend(storeBackend)
	switch backend {
	case store.CONSUL:
		return []string{consulHTTPAddr}
	case store.BOLTDB:
		return []string{boldDBFilePath}
	default:
		return []string{}
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("store-backend", string(store.BOLTDB), "store backend to use boltdb/consul")
	serverCmd.Flags().StringVar(&bindAddress, "bind-address", "0.0.0.0", "bind address for api server")
	serverCmd.Flags().StringVar(&httpPort, "http-port", "1410", "HTTP API port to listen on")
	serverCmd.Flags().StringVar(&storeBasePath, "store-base-path", "dwarka", "Base path for persisting all data")
	serverCmd.Flags().StringVar(&bucketName, "bucket-name", "dwarka", "Base path for persisting all data")

	backend := store.Backend(storeBackend)
	switch backend {
	case store.CONSUL:
		configureAndAddConsulBackendFlags()
	case store.BOLTDB:
		configureAndAddBoltDBFlags()
	}
}

func configureAndAddConsulBackendFlags() {
	consul.Register()
	usage := `The 'address' and port of the Consul HTTP agent. The value can be
an IP address or DNS address, but it must also include the port.`
	serverCmd.Flags().StringVar(&consulHTTPAddr, "consul-http-addr", "http://127.0.0.1:8500", usage)
}

func configureAndAddBoltDBFlags() {
	boltdb.Register()
	serverCmd.Flags().StringVar(&boldDBFilePath, "boltdb-file-path", "data/dwarka", "file path to use for persisting into disk")
}
