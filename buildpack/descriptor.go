// Buildpack descriptor file (https://github.com/buildpacks/spec/blob/main/buildpack.md#buildpacktoml-toml).

package buildpack

import "github.com/BurntSushi/toml"

type Descriptor struct {
	API       string `toml:"api"`
	Buildpack Info   `toml:"buildpack"` // exactly one of 'buildpack' or 'extension' must be populated
	Extension Info   `toml:"extension"` // exactly one of 'buildpack' or 'extension' must be populated
	Order     Order  `toml:"order"`
	Dir       string `toml:"-"`
}

func (b *Descriptor) ConfigFile() *Descriptor {
	return b
}

func (b *Descriptor) IsBuildpack() bool {
	return b.Buildpack.ID != ""
}

func (b *Descriptor) IsExtension() bool {
	return b.Extension.ID != ""
}

func (b *Descriptor) IsMetaBuildpack() bool {
	return len(b.Order) > 0
}

func (b *Descriptor) String() string {
	return b.Buildpack.Name + " " + b.Buildpack.Version
}

type Info struct {
	ClearEnv bool     `toml:"clear-env,omitempty"`
	Homepage string   `toml:"homepage,omitempty"`
	ID       string   `toml:"id"`
	Name     string   `toml:"name"`
	Version  string   `toml:"version"`
	SBOM     []string `toml:"sbom-formats,omitempty" json:"sbom-formats,omitempty"`
}

type Order []Group

type Group struct {
	Group []GroupBuildpack `toml:"group"`
}

func ReadGroup(path string) (Group, error) {
	var group Group
	_, err := toml.DecodeFile(path, &group)
	return group, err
}

func ReadOrder(path string) (Order, Order, error) {
	var order struct {
		Order    Order `toml:"order"`
		OrderExt Order `toml:"order-ext"`
	}
	_, err := toml.DecodeFile(path, &order)
	return order.Order, order.OrderExt, err
}

func (bg Group) Append(group ...Group) Group {
	for _, g := range group {
		bg.Group = append(bg.Group, g.Group...)
	}
	return bg
}

// A GroupBuildpack represents a buildpack referenced in a buildpack.toml's [[order.group]].
// It may be a regular buildpack, or a meta buildpack.
type GroupBuildpack struct {
	API      string `toml:"api,omitempty" json:"-"`
	Homepage string `toml:"homepage,omitempty" json:"homepage,omitempty"`
	ID       string `toml:"id" json:"id"`
	Optional bool   `toml:"optional,omitempty" json:"optional,omitempty"`
	Version  string `toml:"version" json:"version"`
}

func (bp GroupBuildpack) String() string {
	return bp.ID + "@" + bp.Version
}

func (bp GroupBuildpack) NoOpt() GroupBuildpack {
	bp.Optional = false
	return bp
}

func (bp GroupBuildpack) NoAPI() GroupBuildpack {
	bp.API = ""
	return bp
}

func (bp GroupBuildpack) NoHomepage() GroupBuildpack {
	bp.Homepage = ""
	return bp
}
