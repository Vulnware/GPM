package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

func attachDockerContainer(d Docker, gEnv *[]string) {
	log.Printf("d.Compose: %v", d.Compose)
	if d.Compose {
		log.Printf("Attaching to %s in %s with command '%s'\n", d.Name, d.Image, d.Command)
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		labelFilters := filters.NewArgs()
		labelFilters.Add("label", "com.docker.compose.project="+d.Name)

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
			All:     true,
			Filters: labelFilters,
		})
		if err != nil {
			panic(err)
		}
		// log the containers
		log.Printf("Found %d containers for %s\n", len(containers), d.Name)
		for _, container := range containers {
			log.Printf("Attaching to %s in %s with command '%s'\n", container.Names[0], d.Image, d.Command)

			go func(container types.Container) {

				// attach to the container
				attach, err := cli.ContainerAttach(ctx, container.ID, types.ContainerAttachOptions{
					Stream: true,
					Stdout: true,
					Stderr: true,
				})
				if err != nil {
					panic(err)
				}
				defer attach.Close()
				// Define the options to retrieve the container logs
				// defince since 1 day ago
				since := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
				options := types.ContainerLogsOptions{
					ShowStdout: true,
					ShowStderr: true,
					Since:      since,
				}

				// Retrieve the container logs as a stream
				ctx := context.Background()
				reader, err := cli.ContainerLogs(ctx, container.ID, options)
				if err != nil {
					panic(err)
				}

				// Read the logs from the stream and print them to the console
				var buf bytes.Buffer
				_, err = buf.ReadFrom(reader)
				if err != nil {
					panic(err)
				}
				// send the logs to the client
				data := map[string]interface{}{
					"method":  "service",
					"service": container.Names[0],
					"output":  buf.String(),
				}
				appendLogs(container.Names[0], buf.String())
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("%s json marshal error: %s\n", container.Names[0], err)
					return
				}

				// Broadcast message to all clients
				for c := range clients {
					mut.Lock()
					if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
						log.Printf("%s write error: %s\n", container.Names[0], err)
						mut.Unlock()
						continue
					}
					mut.Unlock()
				}

				// create a io.Reader to collect the output from the command (stderr)
				var output bytes.Buffer

				// create a scanner to read the output
				scanner := bufio.NewScanner(attach.Reader)
				go trackDockerUsage(cli, container.Names[0])
				// read the output
				for scanner.Scan() {
					output := scanner.Text()
					// if err := decoder.Decode(&output); err != nil {
					// 	log.Printf("%s stdout read error: %s\n", p.Name, err)
					// 	break
					// }

					data := map[string]interface{}{
						"method":  "service",
						"service": container.Names[0],
						"output":  output,
					}
					appendLogs(container.Names[0], output)
					jsonData, err := json.Marshal(data)
					if err != nil {
						log.Printf("%s json marshal error: %s\n", container.Names[0], err)
						break
					}

					// Broadcast message to all clients
					for c := range clients {
						mut.Lock()
						if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
							log.Printf("%s write error: %s\n", container.Names[0], err)
							mut.Unlock()
							continue
						}
						mut.Unlock()
					}

				}
				if output.Len() > 0 {
					mut.Lock()
					for c := range clients {
						if err := c.WriteMessage(websocket.TextMessage, output.Bytes()); err != nil {
							log.Printf("%s write error: %s\n", container.Names[0], err)
							mut.Unlock()
							continue
						}
					}
					mut.Unlock()
				}
			}(container)
		}

	} else {
		log.Printf("Attaching to %s in %s with command '%s'\n", d.Name, d.Image, d.Command)

		// use docker client to start container
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		// attach to the container
		attach, err := cli.ContainerAttach(ctx, d.Name, types.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		})
		if err != nil {
			panic(err)
		}
		defer attach.Close()
		// read all the logs at once
		// Define the options to retrieve the container logs
		// defince since 1 day ago
		since := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
		options := types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Since:      since,
		}

		// Retrieve the container logs as a stream
		reader, err := cli.ContainerLogs(ctx, d.Name, options)
		if err != nil {
			panic(err)
		}

		// Read the logs from the stream and print them to the console
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		if err != nil {
			panic(err)
		}
		// send the logs to the client
		data := map[string]interface{}{
			"method":  "service",
			"service": d.Name,
			"output":  buf.String(),
		}
		appendLogs(d.Name, buf.String())
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("%s json marshal error: %s\n", d.Name, err)
			return
		}

		// Broadcast message to all clients
		for c := range clients {
			mut.Lock()
			if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Printf("%s write error: %s\n", d.Name, err)
				mut.Unlock()
				continue
			}
			mut.Unlock()
		}
		// create a io.Reader to collect the output from the command (stderr)
		var output bytes.Buffer

		// create a scanner to read the output
		scanner := bufio.NewScanner(attach.Reader)
		go trackDockerUsage(cli, d.Name)
		// read the output
		for scanner.Scan() {
			output := scanner.Text()
			data := map[string]interface{}{
				"method":  "service",
				"service": d.Name,
				"output":  output,
			}
			appendLogs(d.Name, output)
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("%s json marshal error: %s \n", d.Name, err)
				break
			}

			// Broadcast message to all clients
			for c := range clients {
				mut.Lock()
				if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					log.Printf("%s write error: %s \n", d.Name, err)
					mut.Unlock()
					continue
				}
				mut.Unlock()
			}
		}
		if output.Len() > 0 {
			mut.Lock()
			for c := range clients {
				if err := c.WriteMessage(websocket.TextMessage, output.Bytes()); err != nil {
					log.Printf("%s write error: %s \n", d.Name, err)
					mut.Unlock() // prevent deadlock
					continue
				}
			}
			mut.Unlock()
		}
	}
}

func trackDockerUsage(cli *client.Client, name string) {
	// Retrieve the container stats as a stream
	ctx := context.Background()
	reader, err := cli.ContainerStats(ctx, name, true)
	if err != nil {
		panic(err)
	}
	defer reader.Body.Close()

	// Read the container stats from the stream and print them to the console
	var stats types.StatsJSON
	var prevCPU, prevSystem uint64

	for {
		// Decode the stats
		if err := json.NewDecoder(reader.Body).Decode(&stats); err != nil {
			log.Printf("Error decoding stats: %s %s %s %s", err, stats.Read, stats.PreRead, stats.BlkioStats)
		}

		// Calculate CPU percentage
		cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - prevCPU)
		systemDelta := float64(stats.CPUStats.SystemUsage - prevSystem)
		cpuPercent := 0.0
		if systemDelta > 0.0 && cpuDelta > 0.0 {
			cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
		}

		// Print the CPU usage and memory usage of the container
		data := map[string]interface{}{
			"method":  "usage",
			"service": name,
			"cpu":     cpuPercent,
			"memory":  float64(stats.MemoryStats.Usage) / 1024 / 1024,
		}

		// Save current values for next iteration
		prevCPU = stats.CPUStats.CPUUsage.TotalUsage
		prevSystem = stats.CPUStats.SystemUsage

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("%s json marshal error: %s\n", name, err)
			break
		}

		// Broadcast message to all clients
		for c := range clients {
			sendDataToClient(c, jsonData)
		}

		time.Sleep(1 * time.Second)

	}
}
