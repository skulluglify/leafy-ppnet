package app

import "PapayaNet/papaya"

func App(pn papaya.PapayaNetImpl) error {

	if err := pn.Serve("127.0.0.1", 8000); err != nil {

		pn.GetConsole().Error(err)
	}
	return nil
}
