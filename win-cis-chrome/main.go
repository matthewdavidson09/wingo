package main

import (
	"fmt"
	"log"
	"golang.org/x/sys/windows/registry"
)

// Function to create a registry key if it doesn't exist
func createRegistryKeyIfNotExist(path string) error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create registry key: %w", err)
	}
	defer key.Close()
	return nil
}

// Function to set registry values
func setRegistryKey(path string, values map[string]interface{}) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	for name, value := range values {
		switch v := value.(type) {
		case uint32:
			err = key.SetDWordValue(name, v)
		case string:
			err = key.SetStringValue(name, v)
		default:
			log.Printf("unsupported type for key %s", name)
		}
		if err != nil {
			return fmt.Errorf("failed to set registry value %s: %w", name, err)
		}
	}

	return nil
}

func deleteRegistryKey(path, name string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	// Check if the value exists before attempting to delete it
	_, _, err = key.GetValue(name, nil)
	if err == registry.ErrNotExist {
		log.Printf("Registry value %s does not exist, skipping deletion.", name)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check registry value %s: %w", name, err)
	}

	// If the value exists, attempt to delete it
	if err := key.DeleteValue(name); err != nil {
		return fmt.Errorf("failed to delete registry value %s: %w", name, err)
	}

	return nil
}

// Function to check if a registry key exists
func registryKeyExists(path string) bool {
	_, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	return err == nil
}

func main() {
	chromePath := `SOFTWARE\Policies\Google\Chrome`

	// Check if the registry key exists, if not create it
	if !registryKeyExists(chromePath) {
		log.Println("Chrome registry key not found, creating key...")
		err := createRegistryKeyIfNotExist(chromePath)
		if err != nil {
			log.Fatalf("Error creating registry key: %v", err)
		}
	}

	// Set the registry values
	values := map[string]interface{}{
		"SafeBrowsingAllowlistDomains": uint32(0),
		"SafeBrowsingProtectionLevel":  uint32(1),
		"MediaRouterCastAllowAllIPs":   uint32(0),
		"BrowserNetworkTimeQueriesEnabled": uint32(1),
		"AudioSandboxEnabled":              uint32(1),
		"BackgroundModeEnabled":            uint32(0),
		"SafeSitesFilterBehavior":          uint32(1),
		"ChromeVariations":                 uint32(0),
		"CertificateTransparencyEnforcementDisabledForLegacyCas": uint32(0),
		"CertificateTransparencyEnforcementDisabledForCas":       uint32(0),
		"CertificateTransparencyEnforcementDisabledForUrls":      uint32(0),
		"SavingBrowserHistoryDisabled":      uint32(0),
		"DNSInterceptionChecksEnabled":      uint32(1),
		"ComponentUpdatesEnabled":           uint32(1),
		"GloballyScopeHTTPAuthCacheEnabled": uint32(0),
		"EnableOnlineRevocationChecks":      uint32(0),
		"RendererCodeIntegrityEnabled":      uint32(1),
		"CommandLineFlagSecurityWarningsEnabled": uint32(1),
		"ThirdPartyBlockingEnabled":             uint32(1),
		"EnterpriseHardwarePlatform":            uint32(0),
		"ForceEphemeralProfiles":                uint32(0),
		"ImportAutofillFormData":                uint32(0),
		"ImportHomepage":                        uint32(0),
		"ImportSearchEngine":                    uint32(0),
		"HSTSPolicyBypassList":                  "",
		"OverrideSecurityRestrictionsOnInsecureOrigin": uint32(0),
		"LookalikeWarningAllowlistDomains":           uint32(0),
		"SuppressUnsupportedOSWarning":               uint32(0),
		"WebRtcLocalIpsAllowedUrls":                  "",
		"DefaultInsecureContentSetting":              uint32(2),
		"DefaultWebBluetoothGuardSetting":            uint32(2),
		"DefaultWebUsbGuardSetting":                  "",
		"BlockExternalExtensions":                    uint32(1),
		"ExtensionAllowedTypes":                      `["extension", "hosted_app", "platform_app", "theme"]`,
		"AllowCrossOriginAuthPrompt":                 uint32(0),
		"AuthSchemes":                                `["ntlm", "negotiate"]`,
		"CloudPrintProxyEnabled":                     uint32(0),
		"RemoteAccessHostAllowRemoteAccessConnections":    uint32(0),
		"RemoteAccessHostAllowUiAccessForRemoteAssistance": uint32(0),
		"RemoteAccessHostRequireCurtain":                  uint32(0),
		"RemoteAccessHostFirewallTraversal":               uint32(0),
		"RemoteAccessHostAllowClientPairing":              uint32(0),
		"RemoteAccessHostAllowRelayedConnection":          uint32(0),
		"DownloadRestrictions":                            uint32(4),
		"DisableSafeBrowsingProceedAnyway":                uint32(1),
		"ChromeCleanupEnabled":                            uint32(1),
		"SitePerProcess":                                  uint32(1),
		"ForceGoogleSafeSearch":                           uint32(1),
		"RelaunchNotification":                            uint32(2),
		"RelaunchNotificationPeriod":                      uint32(86400000),
		"DefaultGeolocationSetting":                       uint32(3),
		"EnableMediaRouter":                               uint32(1),
		"PaymentMethodQueryEnabled":                       uint32(0),
		"BrowserSignin":                                   uint32(2),
		"ChromeCleanupReportingEnabled":                   uint32(0),
		"AlternateErrorPagesEnabled":                      uint32(0),
		"AllowDeletingBrowserHistory":                     uint32(0),
		"NetworkPredictionOptions":                        uint32(2),
		"MetricsReportingEnabled":                         uint32(0),
		"UrlKeyedAnonymizedDataCollectionEnabled":         uint32(0),
		"CloudPrintSubmitEnabled":                         uint32(0),
		"UserFeedbackAllowed":                             uint32(0),
		"AutofillCreditCardEnabled":                       uint32(0),
		"ImportSavedPasswords":                            uint32(0),
		"DiskCacheSize":                                   uint32(250609664),
	}

	err := setRegistryKey(chromePath, values)
	if err != nil {
		log.Fatalf("Error setting registry keys: %v", err)
	}

	// Delete specific registry keys
	err = deleteRegistryKey(chromePath, "DownloadDirectory")
	if err != nil {
		log.Printf("Error deleting registry value: %v", err)
	}

	err = deleteRegistryKey(chromePath, "PromptForDownloadLocation")
	if err != nil {
		log.Printf("Error deleting registry value: %v", err)
	}

	log.Println("Registry updates complete.")
}
