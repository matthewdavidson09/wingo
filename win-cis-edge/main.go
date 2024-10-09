package main

import (
	"fmt"
	"log"
	"golang.org/x/sys/windows/registry"
)

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

func registryKeyExists(path string) bool {
	_, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	return err == nil
}

func main() {
	edgePath := `SOFTWARE\Policies\Microsoft\Edge`
	edgeUpdatePath := `SOFTWARE\Policies\Microsoft\EdgeUpdate`

	if registryKeyExists(edgePath) {
		values := map[string]interface{}{
			"EnableMediaRouter":                      uint32(0),
			"SpotlightExperiencesAndRecommendationsEnabled": uint32(0),
			"DefaultinsecurecontentSetting":          uint32(2),
			"FeatureFlagoverridesControl":            uint32(0),
			"BasicAuthOverHttpEnabled":               uint32(0),
			"AllowCrossOriginAuthPrompt":             uint32(0),
			"PasswordManagerEnabled":                 uint32(0),
			"StartupBoostEnabled":                    uint32(0),
			"InsecurePrivateNetworkRequestsAllowed":  uint32(0),
			"SmartScreenEnabled":                     uint32(1),
			"SmartScreenPuaEnabled":                  uint32(1),
			"SmartScreenForTrustedDownloadsEnabled":  uint32(1),
			"PreventSmartScreenPromptOverride":       uint32(1),
			"TyposquattingCheckerEnabled":            uint32(1),
			"AdsSettingForIntrusiveAdsSites":         uint32(2),
			"DownloadRestrictions":                   uint32(2),
			"MediaRouterCastAllowAllIPs":             uint32(0),
			"ImportAutofillFormData":                 uint32(0),
			"ImportBrowserSettings":                  uint32(0),
			"ImportHomepage":                         uint32(0),
			"ImportPaymentInfo":                      uint32(0),
			"ImportSavedPasswords":                   uint32(0),
			"ImportSearchEngine":                     uint32(0),
			"EnterpriseHardwarePlatformAPIEnabled":   uint32(0),
			"PersonalizationReportingEnabled":        uint32(0),
			"BrowserNetworkTimeQueriesEnabled":       uint32(1),
			"RemoteDebuggingAllowed":                 uint32(0),
			"LocalProvidersEnabled":                  uint32(0),
			"AudioSandboxEnabled":                    uint32(1),
			"UserFeedbackAllowed":                    uint32(0),
			"PaymentMethodQueryEnabled":              uint32(0),
			"AutoImportAtFirstRun":                   uint32(4),
			"TrackingPrevention":                     uint32(2),
			"ClearBrowsingDataOnExit":                uint32(0),
			"ClearCachedImagesAndFilesOnExit":        uint32(0),
			"HSTSPolicyBypassList":                   "",
			"ConfigureShare":                         uint32(1),
			"InternetExplorerIntegrationComplexNavDataTypes": uint32(0),
			"BackgroundModeEnabled":                  uint32(0),
			"ExperimentationAndConfigurationServiceControl": uint32(0),
			"DeleteDataOnMigration":                  uint32(0),
			"SavingBrowserHistoryDisabled":           uint32(0),
			"SyncDisabled":                           uint32(0),
			"DNSInterceptionChecksEnabled":           uint32(1),
			"AutofillCreditCardEnabled":              uint32(0),
			"BrowserLegacyExtensionPointsBlockingEnabled": uint32(1),
			"ComponentUpdatesEnabled":                uint32(1),
			"AllowDeletingBrowserHistory":            uint32(0),
			"EdgeFollowEnabled":                      uint32(0),
			"GloballyScopeHTTPAuthCacheEnabled":      uint32(0),
			"BrowserGuestModeEnabled":                uint32(0),
			"NetworkPredictionOptions":               uint32(2),
			"BrowserAddProfileEnabled":               uint32(0),
			"RendererCodeIntegrityEnabled":           uint32(1),
			"ResolveNavigationErrorsUseWebService":   uint32(0),
			"CommandLineFlagSecurityWarningsEnabled": uint32(1),
			"SitePerProcess":                         uint32(1),
			"TravelAssistanceEnabled":                uint32(0),
			"ForceEphemeralProfiles":                 uint32(0),
			"InsecureFormsWarningsEnabled":           uint32(1),
			"ForceBingSafeSearch":                    uint32(1),
			"ForceGoogleSafeSearch":                  uint32(1),
			"EnhanceSecurityMode":                    uint32(1),
			"HideFirstRunExperience":                 uint32(1),
			"InAppSupportEnabled":                    uint32(0),
			"RelaunchNotification":                   uint32(2),
			"WebRtcLocalhostIpHandling":              "default_public_interface_only",
			"DiskCacheSize":                          uint32(250000000),
			"RelaunchNotificationPeriod":             uint32(86400000),
			"EdgeShoppingAssistantEnabled":           uint32(0),
			"ExternalProtocolDialogShowAlwaysOpenCheckbox": uint32(0),
			"ShowMicrosoftRewards":                   uint32(0),
			"SharedArrayBufferUnrestrictedAccessAllowed": uint32(0),
			"AlternateErrorPagesEnabled":             uint32(0),
			"SuppressUnsupportedOSWarning":           uint32(0),
		}

		err := setRegistryKey(edgePath, values)
		if err != nil {
			log.Fatalf("Error setting Edge registry keys: %v", err)
		}
	} else {
		log.Println("Edge registry key not found, skipping...")
	}

	// Handle EdgeUpdate registry key
	if registryKeyExists(edgeUpdatePath) {
		updateValues := map[string]interface{}{
			"UpdateDefault": uint32(1),
		}

		err := setRegistryKey(edgeUpdatePath, updateValues)
		if err != nil {
			log.Fatalf("Error setting EdgeUpdate registry key: %v", err)
		}
	} else {
		log.Println("EdgeUpdate registry key not found, skipping...")
	}

	log.Println("Registry updates for Microsoft Edge complete.")
}
