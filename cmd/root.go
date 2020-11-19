/*
Copyright © 2020 weapons97@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/weapons97/helm-api/pkg/helm"
	"github.com/weapons97/helm-api/pkg/middle"
	protos "github.com/weapons97/helm-api/pkg/protos"
	"github.com/weapons97/helm-api/pkg/services"
	"github.com/weapons97/helm-api/pkg/utils"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "helm-api",
	Short: "server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		lis, e := net.Listen("tcp", fmt.Sprintf(`:%s`, viper.GetString(`PORT`)))
		if e != nil {
			utils.Panic().Err(e).Send()
		}
		gser := grpc.NewServer(
			grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					middle.UnaryConnectionWaiter,
					middle.UnaryHelmClient,
				),
			),
		)
		hs := &services.HelmApiService{}
		protos.RegisterHelmApiServiceServer(gser, hs)
		go func() {
			utils.Info().Msgf("listen on: %v of grpc", viper.GetString(`PORT`))
			e = gser.Serve(lis)
			if e != nil {
				utils.Error().Err(e).Send()
			}
			cancel()
		}()
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		e = protos.RegisterHelmApiServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", viper.GetString(`PORT`)), opts)
		if e != nil {
			utils.Panic().Err(e).Send()
		}
		utils.Info().Msgf("listen on: %v of http", viper.GetString(`HTTP_PORT`))
		e = http.ListenAndServe(fmt.Sprintf(`:%s`, viper.GetString(`HTTP_PORT`)), mux)
		if e != nil {
			utils.Panic().Err(e).Send()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if e := rootCmd.Execute(); e != nil {
		utils.Error().Err(e).Send()
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		initConfig,
		initLog,
		initGC,
	)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helm-api.yaml)")
}
func initGC() {
	tmp := viper.GetString(`TMP`)
	e := os.MkdirAll(tmp, 0777)
	if e != nil {
		utils.Error().Str(`dir`, tmp).Err(e).Send()
		os.Exit(1)
	}
	go helm.ContextGC()
}

// initLog parse log level from env and set global log level
func initLog() {
	if viper.GetString(`DEBUG`) == `true` {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}

	l := viper.GetString(`LOGLEVEL`)
	if l == "" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		utils.Warn().Msg(`can't find loglevel'`)
		return
	}

	lv, e := zerolog.ParseLevel(l)
	if e != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		utils.Warn().Interface(`level`, l).Msg(`bad loglevel`)
		return
	}
	zerolog.SetGlobalLevel(lv)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetEnvPrefix(`HELM_API`)

	e := viper.BindEnv(
		`DEBUG`,
	)
	if e != nil {
		panic(e)
	}

	e = viper.BindEnv(
		`LOGLEVEL`,
	)
	if e != nil {
		panic(e)
	}
	viper.SetDefault(`LOGLEVEL`, `info`)

	e = viper.BindEnv(
		`TMP`,
	)
	if e != nil {
		panic(e)
	}
	viper.SetDefault(`TMP`, `/var/tmp/helm-api`)

	e = viper.BindEnv(`PORT`)
	viper.SetDefault(`PORT`, `8848`)
	if e != nil {
		panic(e)
	}

	e = viper.BindEnv(`HTTP_PORT`)
	viper.SetDefault(`HTTP_PORT`, `8611`)
	if e != nil {
		panic(e)
	}

	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		utils.Error().Err(err).Send()
	}
}
