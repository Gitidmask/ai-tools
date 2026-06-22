import React from 'react';

interface Project {
  id: string;
  name: string;
  path: string;
  pinned: boolean;
  last_used: number;
}

export default function WorkspacePanel() {
  const [projects, setProjects] = React.useState<Project[]>([]);

  React.useEffect(() => {
    // Load projects via Wails binding
    async function loadProjects() {
      try {
        const list = await window.runtime?.call('project.ListProjects') as Project[] | undefined;
        if (list) setProjects(list);
      } catch (e) {
        console.error('Failed to load projects', e);
      }
    }
    loadProjects();
  }, []);

  return (
    <div className="workspace-panel">
      <div className="workspace-header">
        <h2>工作区</h2>
      </div>
      <div className="workspace-project-list">
        {projects.map((project) => (
          <div key={project.id} className="project-item">
            <span className="project-name">{project.name}</span>
            <span className="project-path">{project.path}</span>
          </div>
        ))}
        {projects.length === 0 && (
          <p className="workspace-empty">暂无项目</p>
        )}
      </div>
    </div>
  );
}
