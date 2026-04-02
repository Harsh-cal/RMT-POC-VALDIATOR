export const mockReleases = [
  {
    id: "r1",
    label: "R1 — Multiple Issues (FAIL)",
    payload: {
      release_name: "Release_v2.0.0",
      version: "2.0.0",
      target_fleet: "Boeing 737",
      aircraft: {
        tailNumber: "VT-AX1",
        type: "787-9",
        system: "IFE",
        currentSoftware: {
          IFE_Software: {
            version: "2.5.0",
            partNumber: "PN-IFE-250"
          },
          Navigation_Module: {
            version: "2.2.0",
            partNumber: "PN-NAV-220"
          }
        }
      },
      containers: [
        {
          name: "IFE_Software",
          version: "2.0.0",
          partNumber: "PN-IFE-200",
          systemType: "IFE",
          maturity: "beta",
          dependencies: []
        },
        {
          name: "Legacy_Display_Driver",
          version: "1.0.0",
          partNumber: "PN-LEG-100",
          systemType: "IFE",
          maturity: "experimental",
          dependencies: []
        },
        {
          name: "Navigation_Module",
          version: "2.0.0",
          partNumber: "PN-NAV-200",
          systemType: "Connectivity",
          maturity: "stable",
          dependencies: []
        },
        {
          name: "Navigation_Module",
          version: "2.1.0",
          partNumber: "PN-NAV-200",
          systemType: "Connectivity",
          maturity: "stable",
          dependencies: []
        }
      ]
    }
  },
  {
    id: "r2",
    label: "R2 — Safe Release (PASS)",
    payload: {
      release_name: "Release_v3.1.0",
      version: "3.1.0",
      target_fleet: "Boeing 737",
      aircraft: {
        tailNumber: "N7338",
        type: "737-800",
        system: "IFE",
        currentSoftware: {
          IFE_Software: {
            version: "3.0.0",
            partNumber: "PN-IFE-300"
          },
          Navigation_Module: {
            version: "3.5.0",
            partNumber: "PN-NAV-350"
          }
        }
      },
      containers: [
        {
          name: "IFE_Software",
          version: "3.1.0",
          partNumber: "PN-IFE-310",
          systemType: "IFE",
          maturity: "stable",
          dependencies: [
            {
              name: "Navigation_Module",
              required_version: ">=3.5.0"
            }
          ]
        },
        {
          name: "Navigation_Module",
          version: "3.8.0",
          partNumber: "PN-NAV-380",
          systemType: "IFE",
          maturity: "stable",
          dependencies: []
        }
      ]
    }
  },
  {
    id: "r3",
    label: "R3 — Safe Release (PASS)",
    payload: {
      release_name: "Release_v4.0.0",
      version: "4.0.0",
      target_fleet: "Airbus A320",
      aircraft: {
        tailNumber: "G-CYMM",
        type: "A320-200",
        system: "Flight_Control",
        currentSoftware: {
          Flight_Control_UI: {
            version: "3.8.0",
            partNumber: "PN-FC-380"
          },
          Engine_Monitoring_Pack: {
            version: "2.2.0",
            partNumber: "PN-ENG-220"
          }
        }
      },
      containers: [
        {
          name: "Flight_Control_UI",
          version: "4.0.0",
          partNumber: "PN-FC-400",
          systemType: "Flight_Control",
          maturity: "stable",
          dependencies: [
            {
              name: "Engine_Monitoring_Pack",
              required_version: ">=2.0.0"
            }
          ]
        },
        {
          name: "Engine_Monitoring_Pack",
          version: "2.5.0",
          partNumber: "PN-ENG-250",
          systemType: "Flight_Control",
          maturity: "stable",
          dependencies: []
        },
        {
          name: "Diagnostics_Core_V2",
          version: "1.0.0",
          partNumber: "PN-DIAG-100",
          systemType: "Flight_Control",
          maturity: "stable",
          dependencies: []
        }
      ]
    }
  }
];