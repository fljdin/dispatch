Feature: As a Dispatcher user, I want to be able
  to pass a yaml configuration file
  or execute directly through command line

  Scenario: Pass a yaml configuration file to dispatcher
    Given a configuration file in config/TestRunWithConfig.yaml
    When I execute dispatcher
    Then the generated log file in TestRunWithConfig.log should be the same as expected/TestRunWithConfig.log


Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 5 out of 12
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining
