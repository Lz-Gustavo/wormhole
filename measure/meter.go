package measure

import "context"

type Meter struct {
	LatMsr    *LatencyMsr
	StatusMsr *StatusMsr
}

func NewMeter(latfilename, statusfilename string) (*Meter, error) {
	m := &Meter{}
	if latfilename != "" {
		lm, err := NewLatencyMsr(latfilename)
		if err != nil {
			return nil, err
		}
		m.LatMsr = lm
	}

	if statusfilename != "" {
		sm, err := NewStatusMsr(statusfilename)
		if err != nil {
			return nil, err
		}
		m.StatusMsr = sm

		// NOTE: maybe initialize from callee
		go m.StatusMsr.Run(context.TODO())
	}
	return m, nil
}

func (m *Meter) Close() error {
	if m.LatMsr != nil {
		if err := m.LatMsr.Flush(); err != nil {
			return err
		}
		if err := m.LatMsr.Close(); err != nil {
			return err
		}
	}

	if m.StatusMsr != nil {
		if err := m.StatusMsr.Flush(); err != nil {
			return err
		}
		if err := m.StatusMsr.Close(); err != nil {
			return err
		}
	}
	return nil
}
