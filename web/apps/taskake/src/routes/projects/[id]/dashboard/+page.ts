import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { ProjectApi } from '$lib/api/project';
import { DEFAULT_ORG_ID } from '$lib/config';

export const load: PageLoad = async ({ params, fetch }) => {
  const projectId = params.id;
  const projectApi = new ProjectApi(DEFAULT_ORG_ID);
  
  try {
    // Load project data
    const project = await projectApi.getProject(projectId);
    
    // Return data to the page
    return {
      project
    };
  } catch (err) {
    console.error('Error loading project:', err);
    throw error(404, 'Project not found');
  }
};
