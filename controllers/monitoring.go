package controllers

import (
	"errors"
	"fmt"
	"time"
)

type MonitoringController struct {
	BlockcypherController
	BlocktrailController
}

func MakeMonitoringController(blockcypherController BlockcypherController, blocktrailController BlocktrailController) MonitoringController {
	return MonitoringController{
		BlockcypherController: blockcypherController,
		BlocktrailController:  blocktrailController,
	}
}

func (controller MonitoringController) waitForTransfer(address string) (float64, error) {
	ticker := time.NewTicker(3 * time.Minute)
	startingDate := time.Now()
	for {
		select {
		case tick := <-ticker.C:
			fmt.Println(tick.String())
			confirmedBalance, err := controller.BlockcypherController.GetConfirmedBalance(address)

			if err != nil {
				fmt.Println(err)
				confirmedBalance, err = controller.BlocktrailController.GetConfirmedBalance(address)

				fmt.Println(err)
			}

			if confirmedBalance > 0 {
				ticker.Stop()
				return float64(confirmedBalance) / 100000000, nil
			}

			if time.Now().Unix()-startingDate.Unix() == int64(time.Hour.Seconds())*24 {
				ticker.Stop()
				return 0, errors.New("timed out")
			}
		}
	}
}
