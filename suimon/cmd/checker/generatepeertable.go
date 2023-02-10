package checker

import (
	"github.com/bartosian/sui_helpers/suimon/cmd/checker/enums"
	"github.com/bartosian/sui_helpers/suimon/cmd/checker/tablebuilder"
	"github.com/bartosian/sui_helpers/suimon/cmd/checker/tablebuilder/tables"
)

func (checker *Checker) GeneratePeersTable() {
	if !checker.suimonConfig.MonitorsConfig.PeersTable.Display {
		return
	}

	tableConfig := tablebuilder.TableConfig{
		Name:         tables.GetTableTitleSUI(checker.suimonConfig.NetworkType, enums.TableTypePeers),
		Tag:          tables.TableTagSUI,
		Style:        tables.TableStyleSUI,
		RowsCount:    0,
		ColumnsCount: len(tables.ColumnConfigSUI),
		SortConfig:   tables.TableSortConfigSUI,
	}

	columns := make(tablebuilder.Columns, len(tables.ColumnConfigSUI))

	for idx, config := range tables.ColumnConfigSUI {
		columns[idx].Config = config
	}

	for _, peer := range checker.peers {
		if !peer.Metrics.Updated {
			continue
		}

		tableConfig.RowsCount++

		columns[tables.ColumnNameSUIPeerAddress].SetValue(peer.Address)
		columns[tables.ColumnNameSUIPort].SetValue(peer.Port)
		columns[tables.ColumnNameSUITotalTransactions].SetValue(peer.Metrics.TotalTransactionNumber)
		columns[tables.ColumnNameSUIHighestCheckpoints].SetValue(peer.Metrics.HighestSyncedCheckpoint)
		columns[tables.ColumnNameSUIConnectedPeers].SetValue(peer.Metrics.SuiNetworkPeers)
		columns[tables.ColumnNameSUIUptime].SetValue(peer.Metrics.Uptime)
		columns[tables.ColumnNameSUIVersion].SetValue(peer.Metrics.Version)
		columns[tables.ColumnNameSUICommit].SetValue(peer.Metrics.Commit)

		if checker.suimonConfig.HostLookupConfig.EnableLookup {
			columns[tables.ColumnNameSUICountry].SetValue(peer.Location.CountryName)
		} else {
			columns[tables.ColumnNameSUICountry].SetValue(nil)
		}
	}

	if tableConfig.RowsCount == 0 {
		columns.SetNoDataValue()

		tableConfig.RowsCount++
	}

	tableConfig.Columns = columns

	checker.tableBuilderPeer = tablebuilder.NewTableBuilder(tableConfig)
}