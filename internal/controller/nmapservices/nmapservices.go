package nmapservices

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/P1coFly/gRPCnmap/pkg/netvuln_v1"
	"github.com/Ullaakut/nmap"
)

type NmapServices struct {
	log     *slog.Logger
	timeout time.Duration
}

func New(log *slog.Logger, timeout time.Duration) *NmapServices {
	return &NmapServices{log: log, timeout: timeout}
}

func (n *NmapServices) CheckVuln(
	ctx context.Context,
	targets []string,
	ports []int32,
) ([]*netvuln_v1.TargetResult, error) {

	const op = "nmapservices.CheckVuln"

	log := n.log.With(slog.String("op", op))

	log.Info("CheckVuln", "targets", targets, "ports", ports)

	portsStr := make([]string, 0, len(ports))

	for _, p := range ports {
		portsStr = append(portsStr, strconv.Itoa(int(p)))
	}

	gRPCctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithContext(gRPCctx),
		nmap.WithScripts("vulners"),
		nmap.WithTargets(targets...),
		nmap.WithPorts(portsStr...),
		nmap.WithServiceInfo(),
		nmap.WithVersionAll(),
	)

	if err != nil {
		log.Error("unable to create nmap scanner: %v", err)
		return nil, err
	}

	result, warnings, err := scanner.Run()
	if result == nil {
		return nil, fmt.Errorf("Timeout")
	}
	if len(warnings) > 0 {
		log.Warn("run finished with warnings", "warn", warnings) // Warnings are non-critical errors from nmap.
	}
	if err != nil {
		log.Error("unable to run nmap scan: %v", err)
		return nil, err
	}

	// Use the results to print an example output
	var targetResults []*netvuln_v1.TargetResult

	for _, host := range result.Hosts {
		tResult := &netvuln_v1.TargetResult{
			Target: host.Addresses[0].String(),
		}

		for _, port := range host.Ports {

			service := &netvuln_v1.Service{
				Name:    port.Service.Name,
				Version: port.Service.Version,
				TcpPort: int32(port.ID),
			}

			vuk := &netvuln_v1.Vulnerability{}
			for _, scr := range port.Scripts {
				for _, tbl1 := range scr.Tables {
					for _, tbl2 := range tbl1.Tables {

						for _, el := range tbl2.Elements {
							if el.Key == "id" {
								vuk.Identifier = el.Value
							}
							if el.Key == "cvss" {
								v, err := strconv.ParseFloat(el.Value, 32)

								if err != nil {
									log.Error("cann't conv el.val to float32: %v", err)
									continue
								}

								vuk.CvssScore = float32(v)
							}
						}
					}
				}
			}
			if vuk.Identifier != "" {
				service.Vulns = append(service.Vulns, vuk)
			}
			tResult.Services = append(tResult.Services, service)
		}

		targetResults = append(targetResults, tResult)
	}

	return targetResults, nil
}
