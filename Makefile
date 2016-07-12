.PHONY: all
# .SILENT:

all:
	$(MAKE) --no-print-directory -C k8s

clean:
	$(MAKE) --no-print-directory -C k8s clean
