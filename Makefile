.PHONY: all
# .SILENT:

all:
	$(MAKE) --no-print-directory -C cli
	$(MAKE) --no-print-directory -C brokers
	$(MAKE) --no-print-directory -C k8s

test-unit:
	$(MAKE) --no-print-directory -C cli test-unit
	$(MAKE) --no-print-directory -C brokers test-unit
	$(MAKE) --no-print-directory -C k8s test-unit

test:
	$(MAKE) --no-print-directory -C cli test
	$(MAKE) --no-print-directory -C brokers test
	$(MAKE) --no-print-directory -C k8s test

clean:
	$(MAKE) --no-print-directory -C cli clean
	$(MAKE) --no-print-directory -C brokers clean
	$(MAKE) --no-print-directory -C k8s clean
