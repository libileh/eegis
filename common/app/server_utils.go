package app

import (
	"fmt"
	"github.com/libileh/eegis/common/properties"
	"go.uber.org/zap"
	"time"
)

type L10nUtils struct {
	Logger     *zap.SugaredLogger
	Properties *properties.CommonProperties
}

func (l10n *L10nUtils) SetTimeZone() error {

	tz := l10n.Properties.TZ
	loc, err := time.LoadLocation(tz)
	if err != nil {
		l10n.Logger.Warnw("Error loading timezone", "error", err)
		return fmt.Errorf("error on TimeZone %q: %v", tz, err)
	}
	time.Local = loc
	return nil
}
