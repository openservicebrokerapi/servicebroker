.PHONY: all
# .SILENT:

all:
	$(MAKE) --no-print-directory -C k8s

test-unit:
	$(MAKE) --no-print-directory -C k8s test-unit

test:
	$(MAKE) --no-print-directory -C k8s test

clean:
	$(MAKE) --no-print-directory -C k8s clean
