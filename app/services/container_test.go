package services

import (
	"context"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func TestCreateAndRunDockerContainer(t *testing.T) {
	// Создаем тестовый клиент Docker
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Fatalf("Error creating Docker client: %v", err)
	}

	// Подготавливаем контекст
	ctx := context.Background()

	// Тестовая команда для выполнения в контейнере
	testCommand := "echo 'Hello, World!'"

	// Вызываем функцию CreateAndRunDockerContainer
	output, err := CreateAndRunDockerContainer(testCommand)
	if err != nil {
		t.Fatalf("Error running Docker container: %v", err)
	}

	// Проверяем, что вывод содержит ожидаемую строку
	if !strings.Contains(output, "Hello, World!") {
		t.Errorf("Output does not contain expected string")
	}

	// Удаляем контейнер после завершения теста
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		t.Fatalf("Error listing containers: %v", err)
	}
	for _, container := range containers {
		if strings.Contains(container.Image, "ubuntu") {
			if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
				t.Errorf("Error removing container: %v", err)
			}
		}
	}
}
