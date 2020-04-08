package b2b

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccMuleB2bResourceCertificate(t *testing.T) {
	name := "accTest-" + acctest.RandString(5)
	envName := os.Getenv("TEST_ENV_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResourceCertificate_InitialConfig(envName, name),
				Check:  testResourceCertificate_InitialCheck(),
			},
			{
				Config: testResourceCertificate_UpdateConfig(envName, name),
				Check:  testResourceCertificate_UpdateCheck(),
			},
		},
	})
}

func testResourceCertificate_InitialConfig(envName, name string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
}

data "muleb2b_partner" "host" {
  environment_id = data.muleb2b_environment.sbx.id
  host = true
}

data "muleb2b_identifier_type" "as2" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "AS2"
}

resource "muleb2b_partner" "test" {
  name           = "%s"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.as2.id
    value = "%s-id1"
  }
}

resource "muleb2b_certificate" "host" {
  name = "%s-cert-1"
  environment_id = data.muleb2b_environment.sbx.id
  partner_id = data.muleb2b_partner.host.id
  certificate_body = <<EOF
-----BEGIN CERTIFICATE-----
MIIFSjCCAzICCQD0FChujI5RCDANBgkqhkiG9w0BAQsFADBnMQswCQYDVQQGEwJ1
czELMAkGA1UECAwCbW4xDTALBgNVBAcMBG1wbHMxDDAKBgNVBAoMA3RzdDEMMAoG
A1UECwwDdHN0MQwwCgYDVQQDDAN0c3QxEjAQBgkqhkiG9w0BCQEWA3N0YTAeFw0y
MDAyMjcxODQyMTFaFw0yMTAyMjYxODQyMTFaMGcxCzAJBgNVBAYTAnVzMQswCQYD
VQQIDAJtbjENMAsGA1UEBwwEbXBsczEMMAoGA1UECgwDdHN0MQwwCgYDVQQLDAN0
c3QxDDAKBgNVBAMMA3RzdDESMBAGCSqGSIb3DQEJARYDc3RhMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEAsUVj2jmsxeVlR4IRfg8WxjOdKMataqIsEzu6
asx4gX/CrCdrGgJjuYiNSL9cxCuKH+cZ0WJev4a9d3DjIos4T4Vtv6clcV1/fmM7
r/pjZY1+6mYcTZExSA/kdIoKS1RGBLqFM6U9fZy7dH4/VUsgnQ2rVe/RShjaokso
MyKOAcy8qa+pzEHQZkhbgJpzhxA5D9eL2xqmC6fRj0cvZdOoXzvlLAaDwekWWLmy
eA54vM7I1r8BqtSOwWaYQDJZvqnOLNYtgHGr9L+o4Ate9rHTt6agYOMguXSmHbuZ
w9UxmPjqBNa26moq0vBTPUkO//8XU3x3pjyRMhBy9+rkGvC85Fv5iUbqjKMUL7E8
gAZ2UTiYyzhYBLaQ291+hpAICxLvTLmnxUIpIVJZFbvtLknwXktKA4Qppi8/zPo1
N8QGL1Ly9NSfO9Xr4RNlYQ6vEdazDuN8CXFe9W9tRrD9Y0XYpkkOvicbrvw4UNRC
hWjfCoa+UzucmfqYRP03TTXO89sFaG+2I6S+P84sD3qQ84XbG42LFjRKaESyfNsq
vosmW2dvS8rZfTK+3HmhCWCcBYrGyCHR3R1FgXC18Dq3qXhWr7/g1waIN7WhtBgP
irC6/t3ftfVsGDafzzmrvQ/y0TFt2KkiIu0cZPcpMlw6CjVJRF6uNRMJCCWsBOGX
fFi/p1cCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEAf3ScwbiZMuBEGBFlvAdwROPe
B13k3A7ou0VClcslxO8FAkCtH9GOO9XU/iMDz1FUBtM59TlADo8tW1cLIfZ1bQ5z
Vs2WOF18eGPkTbx+uJZ9yUhK7q3ha0QaBop6CU/krT4t/OIHSw1YR2MClmEOOKuE
srieswPG0YU4w6U6GdfUv9i3Tzxht+C8ETUAaxFeOkVKPxqx0g8oqKZXK1BwZIau
BNRVOMWDUy81EkfdsZ63UwAJcUkxSfkdyoMMKcFIkz2bAozbi8xSCw6PQ02T+FYG
PgNNPcr5h4KuSLRVP5rTHiP0zVYkJ56LBaQ8V7EEcRquUCT+E7T/JTJw5dKR2bt3
GKjq93O/rMH//W+JSMDFd1YiqWPCe+hPxJtApA+z/39E31OmnK29yIKMTteRAEOV
GRB+S387t0DcmFrZkg81Fufwg9Zx5ZyYkar/ePR3h9lZEaXOM9VFVlC9UTrfdF84
sQjP3jyeddplRXqADT8x+4iRCjaWo5cQZBu+WPjeY9sSIuscfsk+PS0eG+5+LekK
AnQpD6g2QPpGY2x/3moAnS9IAnxhi3SxxwfLdOrhKMIoIuZaWex/Tbp79U4wblSa
HTUSfNOJpiUGKNGZ4wVNoa6Qki5fH1pWlGE+HTcizBYW/Sv5mHgNIbjOTiPsVvZd
axlm6QIk1IMrzP7+XQo=
-----END CERTIFICATE-----
EOF
}

resource "muleb2b_certificate" "test" {
  name = "%s-cert-2"
  environment_id = data.muleb2b_environment.sbx.id
  partner_id = muleb2b_partner.test.id
  certificate_body = <<EOF
-----BEGIN CERTIFICATE-----
MIIFSjCCAzICCQD0FChujI5RCDANBgkqhkiG9w0BAQsFADBnMQswCQYDVQQGEwJ1
czELMAkGA1UECAwCbW4xDTALBgNVBAcMBG1wbHMxDDAKBgNVBAoMA3RzdDEMMAoG
A1UECwwDdHN0MQwwCgYDVQQDDAN0c3QxEjAQBgkqhkiG9w0BCQEWA3N0YTAeFw0y
MDAyMjcxODQyMTFaFw0yMTAyMjYxODQyMTFaMGcxCzAJBgNVBAYTAnVzMQswCQYD
VQQIDAJtbjENMAsGA1UEBwwEbXBsczEMMAoGA1UECgwDdHN0MQwwCgYDVQQLDAN0
c3QxDDAKBgNVBAMMA3RzdDESMBAGCSqGSIb3DQEJARYDc3RhMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEAsUVj2jmsxeVlR4IRfg8WxjOdKMataqIsEzu6
asx4gX/CrCdrGgJjuYiNSL9cxCuKH+cZ0WJev4a9d3DjIos4T4Vtv6clcV1/fmM7
r/pjZY1+6mYcTZExSA/kdIoKS1RGBLqFM6U9fZy7dH4/VUsgnQ2rVe/RShjaokso
MyKOAcy8qa+pzEHQZkhbgJpzhxA5D9eL2xqmC6fRj0cvZdOoXzvlLAaDwekWWLmy
eA54vM7I1r8BqtSOwWaYQDJZvqnOLNYtgHGr9L+o4Ate9rHTt6agYOMguXSmHbuZ
w9UxmPjqBNa26moq0vBTPUkO//8XU3x3pjyRMhBy9+rkGvC85Fv5iUbqjKMUL7E8
gAZ2UTiYyzhYBLaQ291+hpAICxLvTLmnxUIpIVJZFbvtLknwXktKA4Qppi8/zPo1
N8QGL1Ly9NSfO9Xr4RNlYQ6vEdazDuN8CXFe9W9tRrD9Y0XYpkkOvicbrvw4UNRC
hWjfCoa+UzucmfqYRP03TTXO89sFaG+2I6S+P84sD3qQ84XbG42LFjRKaESyfNsq
vosmW2dvS8rZfTK+3HmhCWCcBYrGyCHR3R1FgXC18Dq3qXhWr7/g1waIN7WhtBgP
irC6/t3ftfVsGDafzzmrvQ/y0TFt2KkiIu0cZPcpMlw6CjVJRF6uNRMJCCWsBOGX
fFi/p1cCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEAf3ScwbiZMuBEGBFlvAdwROPe
B13k3A7ou0VClcslxO8FAkCtH9GOO9XU/iMDz1FUBtM59TlADo8tW1cLIfZ1bQ5z
Vs2WOF18eGPkTbx+uJZ9yUhK7q3ha0QaBop6CU/krT4t/OIHSw1YR2MClmEOOKuE
srieswPG0YU4w6U6GdfUv9i3Tzxht+C8ETUAaxFeOkVKPxqx0g8oqKZXK1BwZIau
BNRVOMWDUy81EkfdsZ63UwAJcUkxSfkdyoMMKcFIkz2bAozbi8xSCw6PQ02T+FYG
PgNNPcr5h4KuSLRVP5rTHiP0zVYkJ56LBaQ8V7EEcRquUCT+E7T/JTJw5dKR2bt3
GKjq93O/rMH//W+JSMDFd1YiqWPCe+hPxJtApA+z/39E31OmnK29yIKMTteRAEOV
GRB+S387t0DcmFrZkg81Fufwg9Zx5ZyYkar/ePR3h9lZEaXOM9VFVlC9UTrfdF84
sQjP3jyeddplRXqADT8x+4iRCjaWo5cQZBu+WPjeY9sSIuscfsk+PS0eG+5+LekK
AnQpD6g2QPpGY2x/3moAnS9IAnxhi3SxxwfLdOrhKMIoIuZaWex/Tbp79U4wblSa
HTUSfNOJpiUGKNGZ4wVNoa6Qki5fH1pWlGE+HTcizBYW/Sv5mHgNIbjOTiPsVvZd
axlm6QIk1IMrzP7+XQo=
-----END CERTIFICATE-----
EOF
}
`, envName, name, name, name, name)
}

func testResourceCertificate_InitialCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		certId := s.Modules[0].Resources["muleb2b_certificate.host"].Primary.ID

		if certId == "" {
			return fmt.Errorf("certificate.Id is empty")
		}

		return nil
	}
}

func testResourceCertificate_UpdateConfig(envName, name string) string {
	return fmt.Sprintf(`data "muleb2b_environment" "sbx" {
  name = "%s"
}

data "muleb2b_partner" "host" {
  environment_id = data.muleb2b_environment.sbx.id
  host = true
}

data "muleb2b_identifier_type" "as2" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "AS2"
}

resource "muleb2b_partner" "test" {
  name           = "%s"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.as2.id
    value = "%s-id1"
  }
}

resource "muleb2b_certificate" "host" {
  name = "%s-cert-1"
  environment_id = data.muleb2b_environment.sbx.id
  partner_id = data.muleb2b_partner.host.id
  certificate_body = <<EOF
-----BEGIN CERTIFICATE-----
MIIFSjCCAzICCQD0FChujI5RCDANBgkqhkiG9w0BAQsFADBnMQswCQYDVQQGEwJ1
czELMAkGA1UECAwCbW4xDTALBgNVBAcMBG1wbHMxDDAKBgNVBAoMA3RzdDEMMAoG
A1UECwwDdHN0MQwwCgYDVQQDDAN0c3QxEjAQBgkqhkiG9w0BCQEWA3N0YTAeFw0y
MDAyMjcxODQyMTFaFw0yMTAyMjYxODQyMTFaMGcxCzAJBgNVBAYTAnVzMQswCQYD
VQQIDAJtbjENMAsGA1UEBwwEbXBsczEMMAoGA1UECgwDdHN0MQwwCgYDVQQLDAN0
c3QxDDAKBgNVBAMMA3RzdDESMBAGCSqGSIb3DQEJARYDc3RhMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEAsUVj2jmsxeVlR4IRfg8WxjOdKMataqIsEzu6
asx4gX/CrCdrGgJjuYiNSL9cxCuKH+cZ0WJev4a9d3DjIos4T4Vtv6clcV1/fmM7
r/pjZY1+6mYcTZExSA/kdIoKS1RGBLqFM6U9fZy7dH4/VUsgnQ2rVe/RShjaokso
MyKOAcy8qa+pzEHQZkhbgJpzhxA5D9eL2xqmC6fRj0cvZdOoXzvlLAaDwekWWLmy
eA54vM7I1r8BqtSOwWaYQDJZvqnOLNYtgHGr9L+o4Ate9rHTt6agYOMguXSmHbuZ
w9UxmPjqBNa26moq0vBTPUkO//8XU3x3pjyRMhBy9+rkGvC85Fv5iUbqjKMUL7E8
gAZ2UTiYyzhYBLaQ291+hpAICxLvTLmnxUIpIVJZFbvtLknwXktKA4Qppi8/zPo1
N8QGL1Ly9NSfO9Xr4RNlYQ6vEdazDuN8CXFe9W9tRrD9Y0XYpkkOvicbrvw4UNRC
hWjfCoa+UzucmfqYRP03TTXO89sFaG+2I6S+P84sD3qQ84XbG42LFjRKaESyfNsq
vosmW2dvS8rZfTK+3HmhCWCcBYrGyCHR3R1FgXC18Dq3qXhWr7/g1waIN7WhtBgP
irC6/t3ftfVsGDafzzmrvQ/y0TFt2KkiIu0cZPcpMlw6CjVJRF6uNRMJCCWsBOGX
fFi/p1cCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEAf3ScwbiZMuBEGBFlvAdwROPe
B13k3A7ou0VClcslxO8FAkCtH9GOO9XU/iMDz1FUBtM59TlADo8tW1cLIfZ1bQ5z
Vs2WOF18eGPkTbx+uJZ9yUhK7q3ha0QaBop6CU/krT4t/OIHSw1YR2MClmEOOKuE
srieswPG0YU4w6U6GdfUv9i3Tzxht+C8ETUAaxFeOkVKPxqx0g8oqKZXK1BwZIau
BNRVOMWDUy81EkfdsZ63UwAJcUkxSfkdyoMMKcFIkz2bAozbi8xSCw6PQ02T+FYG
PgNNPcr5h4KuSLRVP5rTHiP0zVYkJ56LBaQ8V7EEcRquUCT+E7T/JTJw5dKR2bt3
GKjq93O/rMH//W+JSMDFd1YiqWPCe+hPxJtApA+z/39E31OmnK29yIKMTteRAEOV
GRB+S387t0DcmFrZkg81Fufwg9Zx5ZyYkar/ePR3h9lZEaXOM9VFVlC9UTrfdF84
sQjP3jyeddplRXqADT8x+4iRCjaWo5cQZBu+WPjeY9sSIuscfsk+PS0eG+5+LekK
AnQpD6g2QPpGY2x/3moAnS9IAnxhi3SxxwfLdOrhKMIoIuZaWex/Tbp79U4wblSa
HTUSfNOJpiUGKNGZ4wVNoa6Qki5fH1pWlGE+HTcizBYW/Sv5mHgNIbjOTiPsVvZd
axlm6QIk1IMrzP7+XQo=
-----END CERTIFICATE-----
EOF
}

resource "muleb2b_certificate" "test" {
  name = "%s-cert-3"
  environment_id = data.muleb2b_environment.sbx.id
  partner_id = muleb2b_partner.test.id
  certificate_body = <<EOF
-----BEGIN CERTIFICATE-----
MIIFSjCCAzICCQD0FChujI5RCDANBgkqhkiG9w0BAQsFADBnMQswCQYDVQQGEwJ1
czELMAkGA1UECAwCbW4xDTALBgNVBAcMBG1wbHMxDDAKBgNVBAoMA3RzdDEMMAoG
A1UECwwDdHN0MQwwCgYDVQQDDAN0c3QxEjAQBgkqhkiG9w0BCQEWA3N0YTAeFw0y
MDAyMjcxODQyMTFaFw0yMTAyMjYxODQyMTFaMGcxCzAJBgNVBAYTAnVzMQswCQYD
VQQIDAJtbjENMAsGA1UEBwwEbXBsczEMMAoGA1UECgwDdHN0MQwwCgYDVQQLDAN0
c3QxDDAKBgNVBAMMA3RzdDESMBAGCSqGSIb3DQEJARYDc3RhMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEAsUVj2jmsxeVlR4IRfg8WxjOdKMataqIsEzu6
asx4gX/CrCdrGgJjuYiNSL9cxCuKH+cZ0WJev4a9d3DjIos4T4Vtv6clcV1/fmM7
r/pjZY1+6mYcTZExSA/kdIoKS1RGBLqFM6U9fZy7dH4/VUsgnQ2rVe/RShjaokso
MyKOAcy8qa+pzEHQZkhbgJpzhxA5D9eL2xqmC6fRj0cvZdOoXzvlLAaDwekWWLmy
eA54vM7I1r8BqtSOwWaYQDJZvqnOLNYtgHGr9L+o4Ate9rHTt6agYOMguXSmHbuZ
w9UxmPjqBNa26moq0vBTPUkO//8XU3x3pjyRMhBy9+rkGvC85Fv5iUbqjKMUL7E8
gAZ2UTiYyzhYBLaQ291+hpAICxLvTLmnxUIpIVJZFbvtLknwXktKA4Qppi8/zPo1
N8QGL1Ly9NSfO9Xr4RNlYQ6vEdazDuN8CXFe9W9tRrD9Y0XYpkkOvicbrvw4UNRC
hWjfCoa+UzucmfqYRP03TTXO89sFaG+2I6S+P84sD3qQ84XbG42LFjRKaESyfNsq
vosmW2dvS8rZfTK+3HmhCWCcBYrGyCHR3R1FgXC18Dq3qXhWr7/g1waIN7WhtBgP
irC6/t3ftfVsGDafzzmrvQ/y0TFt2KkiIu0cZPcpMlw6CjVJRF6uNRMJCCWsBOGX
fFi/p1cCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEAf3ScwbiZMuBEGBFlvAdwROPe
B13k3A7ou0VClcslxO8FAkCtH9GOO9XU/iMDz1FUBtM59TlADo8tW1cLIfZ1bQ5z
Vs2WOF18eGPkTbx+uJZ9yUhK7q3ha0QaBop6CU/krT4t/OIHSw1YR2MClmEOOKuE
srieswPG0YU4w6U6GdfUv9i3Tzxht+C8ETUAaxFeOkVKPxqx0g8oqKZXK1BwZIau
BNRVOMWDUy81EkfdsZ63UwAJcUkxSfkdyoMMKcFIkz2bAozbi8xSCw6PQ02T+FYG
PgNNPcr5h4KuSLRVP5rTHiP0zVYkJ56LBaQ8V7EEcRquUCT+E7T/JTJw5dKR2bt3
GKjq93O/rMH//W+JSMDFd1YiqWPCe+hPxJtApA+z/39E31OmnK29yIKMTteRAEOV
GRB+S387t0DcmFrZkg81Fufwg9Zx5ZyYkar/ePR3h9lZEaXOM9VFVlC9UTrfdF84
sQjP3jyeddplRXqADT8x+4iRCjaWo5cQZBu+WPjeY9sSIuscfsk+PS0eG+5+LekK
AnQpD6g2QPpGY2x/3moAnS9IAnxhi3SxxwfLdOrhKMIoIuZaWex/Tbp79U4wblSa
HTUSfNOJpiUGKNGZ4wVNoa6Qki5fH1pWlGE+HTcizBYW/Sv5mHgNIbjOTiPsVvZd
axlm6QIk1IMrzP7+XQo=
-----END CERTIFICATE-----
EOF
}`, envName, name, name, name, name)
}

func testResourceCertificate_UpdateCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		certId := s.Modules[0].Resources["muleb2b_certificate.test"].Primary.ID
		if certId == "" {
			return fmt.Errorf("muleb2b_document_flow.test.config.source_doc_type_id is empty")
		}
		return nil
	}
}
