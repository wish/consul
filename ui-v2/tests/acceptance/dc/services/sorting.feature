@setupApplicationTest
Feature: dc / services / sorting
  Scenario:
    Given 1 datacenter model with the value "dc-1"
    And 12 service models from yaml
    ---
    - Name: Service-A
    - Name: Service-A-proxy
      Kind: 'connect-proxy'
    - Name: Service-B
    - Name: Service-B-proxy
      Kind: 'connect-proxy'
    - Name: Service-C
    - Name: Service-C-proxy
      Kind: 'connect-proxy'
    - Name: Service-D
    - Name: Service-D-proxy
      Kind: 'connect-proxy'
    - Name: Service-E
    - Name: Service-E-proxy
      Kind: 'connect-proxy'
    - Name: Service-F
    - Name: Service-F-proxy
      Kind: 'connect-proxy'
    ---
    When I visit the services page for yaml
    ---
      dc: dc-1
    ---
    When I click selected on the sort
    When I click options.1.button on the sort
    Then I see name on the services vertically like yaml
    ---
    - Service-F
    - Service-E
    - Service-D
    - Service-C
    - Service-B
    - Service-A
    ---
    When I click selected on the sort
    When I click options.0.button on the sort
    Then I see name on the services vertically like yaml
    ---
    - Service-A
    - Service-B
    - Service-C
    - Service-D
    - Service-E
    - Service-F
    ---
