package bakery

// FCOS support tier constants for fedora-sysexts/community extensions.
const (
	// FCOSTierCommunity marks extensions maintained by the community with
	// automated builds from the fedora-sysexts/community repo.
	FCOSTierCommunity = "Community Maintained"
)

// fcosCatalog is the static curated catalog of common fedora-sysexts/community
// extensions. Descriptions are based on the upstream project descriptions and
// the community repo's build configurations.
//
// When extensions not present here are encountered, FCOSLookup returns ok=false
// and the raw GitHub release body is used as fallback.
var fcosCatalog = map[string]ExtensionMeta{
	"docker-ce": {
		Category:    "Container Runtime",
		SupportTier: FCOSTierCommunity,
		Short:       "Docker CE — full Docker engine for Fedora CoreOS",
		Long:        "Ships the complete Docker CE engine (daemon, CLI, containerd, runc) packaged as a sysext for FCOS. Provides the standard Docker experience on an immutable Fedora CoreOS host.",
		Caveats:     nil,
	},
	"tailscale": {
		Category:    "Networking",
		SupportTier: FCOSTierCommunity,
		Short:       "Tailscale — zero-config WireGuard mesh VPN",
		Long:        "Tailscale creates an encrypted WireGuard mesh network between your nodes with no manual key exchange or firewall rules. Ships tailscale and tailscaled for FCOS. Authenticate nodes with tailscale up after provisioning.",
		Caveats:     nil,
	},
	"vscode": {
		Category:    "Development",
		SupportTier: FCOSTierCommunity,
		Short:       "Visual Studio Code — Microsoft's code editor as a sysext",
		Long:        "Ships Visual Studio Code packaged as a system extension for Fedora CoreOS. Enables running VS Code directly on FCOS nodes for development and debugging.",
		Caveats:     nil,
	},
	"vscodium": {
		Category:    "Development",
		SupportTier: FCOSTierCommunity,
		Short:       "VSCodium — free/libre VS Code without Microsoft telemetry",
		Long:        "VSCodium is a community-driven, freely-licensed distribution of VS Code without Microsoft branding, telemetry, or licensing restrictions. Packaged as a sysext for Fedora CoreOS.",
		Caveats:     nil,
	},
	"cilium-cli": {
		Category:    "Networking",
		SupportTier: FCOSTierCommunity,
		Short:       "Cilium CLI — eBPF-based Kubernetes networking and observability",
		Long:        "Cilium provides eBPF-powered networking, security policy enforcement, and observability for Kubernetes clusters. Ships the cilium CLI for FCOS.",
		Caveats:     nil,
	},
	"1password-cli": {
		Category:    "Security",
		SupportTier: FCOSTierCommunity,
		Short:       "1Password CLI — manage secrets and credentials from the terminal",
		Long:        "The 1Password CLI (op) lets you manage passwords, credentials, and secrets from the command line. Integrates with 1Password vaults for secure secret management on FCOS nodes.",
		Caveats:     nil,
	},
	"google-chrome": {
		Category:    "Browser",
		SupportTier: FCOSTierCommunity,
		Short:       "Google Chrome — Chrome browser as a system extension",
		Long:        "Ships Google Chrome packaged as a system extension for Fedora CoreOS.",
		Caveats:     nil,
	},
	"microsoft-edge": {
		Category:    "Browser",
		SupportTier: FCOSTierCommunity,
		Short:       "Microsoft Edge — Edge browser as a system extension",
		Long:        "Ships Microsoft Edge packaged as a system extension for Fedora CoreOS.",
		Caveats:     nil,
	},
	"netbird": {
		Category:    "Networking",
		SupportTier: FCOSTierCommunity,
		Short:       "NetBird — WireGuard-based overlay network with SSO integration",
		Long:        "NetBird creates a peer-to-peer WireGuard overlay network with identity-based access control. Ships the netbird daemon for FCOS. Requires a NetBird management server or the hosted service.",
		Caveats:     nil,
	},
	"netbird-ui": {
		Category:    "Networking",
		SupportTier: FCOSTierCommunity,
		Short:       "NetBird UI — graphical interface for NetBird VPN",
		Long:        "Ships the NetBird graphical tray application for managing NetBird VPN connections on FCOS. Requires the netbird sysext.",
		Caveats:     []string{"Requires the netbird sysext."},
	},
	"nordvpn": {
		Category:    "Networking",
		SupportTier: FCOSTierCommunity,
		Short:       "NordVPN — NordVPN client daemon for FCOS",
		Long:        "Ships the NordVPN client daemon packaged as a system extension for Fedora CoreOS. Requires a NordVPN subscription.",
		Caveats:     nil,
	},
	"nordvpn-gui": {
		Category:    "Networking",
		SupportTier: FCOSTierCommunity,
		Short:       "NordVPN GUI — graphical NordVPN client",
		Long:        "Ships the NordVPN graphical client for Fedora CoreOS. Requires the nordvpn sysext and a NordVPN subscription.",
		Caveats:     []string{"Requires the nordvpn sysext."},
	},
	"bitwarden": {
		Category:    "Security",
		SupportTier: FCOSTierCommunity,
		Short:       "Bitwarden — open-source password manager desktop app",
		Long:        "Ships the Bitwarden desktop application as a system extension for Fedora CoreOS.",
		Caveats:     nil,
	},
	"glab": {
		Category:    "Development",
		SupportTier: FCOSTierCommunity,
		Short:       "GLab — GitLab CLI for managing repositories and CI/CD",
		Long:        "GLab is the official GitLab CLI tool for managing merge requests, issues, pipelines, and repositories from the terminal.",
		Caveats:     nil,
	},
	"virtctl": {
		Category:    "Virtualization",
		SupportTier: FCOSTierCommunity,
		Short:       "virtctl — KubeVirt CLI for managing virtual machines on Kubernetes",
		Long:        "virtctl is the command-line tool for KubeVirt, enabling management of virtual machines running on Kubernetes clusters.",
		Caveats:     nil,
	},
	"cloud-hypervisor": {
		Category:    "Virtualization",
		SupportTier: FCOSTierCommunity,
		Short:       "Cloud Hypervisor — lightweight open-source VMM for cloud workloads",
		Long:        "Cloud Hypervisor is a lightweight virtual machine monitor built on top of KVM, designed for running modern cloud workloads with minimal overhead.",
		Caveats:     nil,
	},
	"openconnect": {
		Category:    "Networking",
		SupportTier: FCOSTierCommunity,
		Short:       "OpenConnect — multi-protocol SSL VPN client",
		Long:        "OpenConnect is a VPN client supporting multiple protocols including AnyConnect, Juniper, GlobalProtect, and Fortinet. Ships as a sysext for Fedora CoreOS.",
		Caveats:     nil,
	},
	"1password-gui": {
		Category:    "Security",
		SupportTier: FCOSTierCommunity,
		Short:       "1Password — desktop password manager application",
		Long:        "Ships the 1Password desktop application as a system extension for Fedora CoreOS.",
		Caveats:     nil,
	},
	"littlesnitch": {
		Category:    "Security",
		SupportTier: FCOSTierCommunity,
		Short:       "Little Snitch — network traffic monitor and firewall",
		Long:        "Little Snitch monitors and controls outgoing network connections, alerting you when applications attempt to connect to remote servers.",
		Caveats:     nil,
	},
}

// FCOSLookup returns curated metadata for an FCOS sysext extension.
// Returns the metadata and true if found; zero ExtensionMeta and false if not.
func FCOSLookup(name string) (ExtensionMeta, bool) {
	meta, ok := fcosCatalog[name]
	return meta, ok
}
