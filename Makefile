.PHONY: all
# .SILENT:

all:
	$(MAKE) --no-print-directory -C k8s
	$(MAKE) --no-print-directory -C cli

test-unit:
	$(MAKE) --no-print-directory -C k8s test-unit
	$(MAKE) --no-print-directory -C cli test-unit

test:
	$(MAKE) --no-print-directory -C k8s test
	$(MAKE) --no-print-directory -C cli test

clean:
	$(MAKE) --no-print-directory -C k8s clean
	$(MAKE) --no-print-directory -C cli clean
