-- Insert seed teams
INSERT INTO teams (uuid, name, description, created_at, updated_at) VALUES
('111e4567-e89b-12d3-a456-426614174000', 'Time de Desenvolvimento', 'Equipe responsável pelo desenvolvimento de features e manutenção do código', '2025-12-01 18:20:00', '2025-12-01 18:20:00'),
('222e4567-e89b-12d3-a456-426614174000', 'Time de DevOps', 'Equipe responsável por infraestrutura, CI/CD e deploy', '2025-12-01 18:20:00', '2025-12-01 18:20:00'),
('333e4567-e89b-12d3-a456-426614174000', 'Time de QA', 'Equipe responsável por testes e garantia de qualidade', '2025-12-01 18:20:00', '2025-12-01 18:20:00'),
('444e4567-e89b-12d3-a456-426614174000', 'Time de UX/UI', 'Equipe responsável por design e experiência do usuário', '2025-12-01 18:20:00', '2025-12-01 18:20:00');


-- Insert seed tasks with various statuses
INSERT INTO tasks (uuid, title, description, status, started_at, finished_at, team_id, created_at, updated_at) VALUES
-- Tasks without team (no team assigned)
('123e4567-e89b-12d3-a456-426614174000', 'Implementar autenticação', 'Criar sistema de autenticação JWT para a API', 'to_do', NULL, NULL, NULL, '2025-12-01 18:21:06', '2025-12-01 18:21:06'),

-- Development Team tasks (team_id = 1)
('123e4567-e89b-12d3-a456-426614174001', 'Criar documentação da API', 'Documentar todos os endpoints da API usando Swagger', 'in_progress', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '2 days', NULL, 1, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '3 days', '2025-12-01 18:21:06'),
('123e4567-e89b-12d3-a456-426614174004', 'Adicionar testes unitários', 'Escrever testes unitários para todas as funções principais', 'to_do', NULL, NULL, 1, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day'),
('123e4567-e89b-12d3-a456-426614174005', 'Otimizar queries do banco', 'Analisar e otimizar queries lentas do banco de dados', 'in_progress', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day', NULL, 1, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '2 days', '2025-12-01 18:21:06'),
('223e4567-e89b-12d3-a456-426614174000', 'Refatorar módulo de autenticação', 'Melhorar estrutura e organização do código de autenticação', 'to_do', NULL, NULL, 1, '2025-12-01 18:21:06', '2025-12-01 18:21:06'),
('223e4567-e89b-12d3-a456-426614174001', 'Implementar feature de notificações', 'Criar sistema de notificações em tempo real', 'in_progress', NULL, NULL, 1, '2025-12-01 19:21:06', '2025-12-01 19:21:06'),

-- DevOps Team tasks (team_id = 2)
('123e4567-e89b-12d3-a456-426614174002', 'Configurar CI/CD', 'Configurar pipeline de CI/CD usando GitHub Actions', 'done', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '5 days', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day', 2, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '7 days', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day'),
('123e4567-e89b-12d3-a456-426614174003', 'Implementar cache Redis', 'Adicionar cache Redis para melhorar performance', 'canceled', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '4 days', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '3 days', 2, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '6 days', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '3 days'),
('123e4567-e89b-12d3-a456-426614174006', 'Criar dashboard de métricas', 'Implementar dashboard para visualizar métricas da aplicação', 'done', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '10 days', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '2 days', 2, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '12 days', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '2 days'),
('323e4567-e89b-12d3-a456-426614174000', 'Configurar monitoramento de logs', 'Implementar sistema centralizado de logs com ELK Stack', 'to_do', NULL, NULL, 2, '2025-12-01 18:21:06', '2025-12-01 18:21:06'),
('323e4567-e89b-12d3-a456-426614174001', 'Otimizar configuração do Docker', 'Melhorar Dockerfile e docker-compose para produção', 'in_progress', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day', NULL, 2, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '2 days', '2025-12-01 18:21:06'),

-- QA Team tasks (team_id = 3)
('423e4567-e89b-12d3-a456-426614174000', 'Criar testes de integração', 'Desenvolver suite completa de testes de integração', 'to_do', NULL, NULL, 3, '2025-12-01 18:21:06', '2025-12-01 18:21:06'),
('423e4567-e89b-12d3-a456-426614174001', 'Executar testes de carga', 'Realizar testes de performance e carga na aplicação', 'in_progress', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day', NULL, 3, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '2 days', '2025-12-01 18:21:06'),
('423e4567-e89b-12d3-a456-426614174002', 'Revisar cobertura de testes', 'Auditar e melhorar cobertura de testes do projeto', 'done', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '3 days', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day', 3, TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '5 days', TIMESTAMP '2025-12-01 18:21:06' - INTERVAL '1 day');
