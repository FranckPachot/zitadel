package setup

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	createAssets = `
CREATE TABLE system.assets (
    instance_id TEXT,
    asset_type TEXT,
    resource_owner TEXT,
    name TEXT,
    content_type TEXT,
    hash TEXT GENERATED ALWAYS AS (md5(data)) STORED,
    data BYTEA,
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (instance_id, resource_owner, name)
);
CREATE OR REPLACE FUNCTION set_asset_hash()
RETURNS TRIGGER AS $$
BEGIN
    NEW.hash := md5(NEW.data);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trigger_set_asset_hash
BEFORE INSERT OR UPDATE ON system.assets
FOR EACH ROW
EXECUTE FUNCTION set_asset_hash();
`
)

type AssetTable struct {
	dbClient *sql.DB
}

func (mig *AssetTable) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, createAssets)
	return err
}

func (mig *AssetTable) String() string {
	return "02_assets"
}
