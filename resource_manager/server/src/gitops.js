import { join, dirname } from "path";
import fs from 'fs';
import simpleGit from 'simple-git';

const localPath = "/tmp/gitops-repo";
const clusterPath = "/clusters/maritime-flow";
const fullPath = join(localPath, clusterPath);

function authUrlWithPat() {
  if (!process.env.GITHUB_PAT) return process.env.GITHUB_REPO;
  if (process.env.GITHUB_REPO.startsWith('https://')) {
    return process.env.GITHUB_REPO.replace('https://', `https://${encodeURIComponent(process.env.GITHUB_PAT)}@`);
  }
  return process.env.GITHUB_REPO;
}

async function ensureRepo() {
  const branch = 'main';
  const git = simpleGit();
  
  if (!fs.existsSync(localPath)) {
    fs.mkdirSync(localPath, { recursive: true });
    const cloneUrl = authUrlWithPat();
    await git.clone(cloneUrl, localPath, ['--branch', branch, '--single-branch']);
  }

  const repoGit = simpleGit({ baseDir: localPath });
  repoGit.addConfig('user.name', 'Maritime Flow GitOps Bot').addConfig('user.email', 'maritimeflowfurb@gmail.com')
  const remoteUrlWithAuth = authUrlWithPat();
  try {
    await repoGit.removeRemote('origin').catch(() => { });
    await repoGit.addRemote('origin', remoteUrlWithAuth).catch(() => { });
  } catch (err) {
    // ignore remote setup errors and continue
  }

  await repoGit.fetch('origin', branch);
  const status = await repoGit.status();
  if (status.current !== branch) {
    await repoGit.checkout(['-B', branch, `origin/${branch}`]).catch(() => { });
  }
  await repoGit.pull('origin', branch).catch(() => { });

  return repoGit;
}

export async function commitManifest({
  directory,
  files,
  replacements
}) {
  const branch = 'main';
  const commitAuthor = { name: 'Maritime Flow GitOps Bot', email: 'maritimeflowfurb@gmail.com' };
  const commitMessage = `Created resource at ${directory}`;
  const repoGit = await ensureRepo();

  for (const [templateFile, outputFile] of Object.entries(files)) {
    const templatePath = join(fullPath, templateFile);
    const outputPath = join(fullPath, outputFile);

    if (!fs.existsSync(templatePath)) {
      throw new Error(`Template not found: ${templateFile}`);
    }

    let content = fs.readFileSync(templatePath, 'utf8');
    for (const [key, value] of Object.entries(replacements)) {
      const regex = new RegExp(`{{\\s*${escapeRegExp(key)}\\s*}}`, "g");
      content = content.replace(regex, value);
    }

    const outDir = dirname(outputPath);
    fs.mkdirSync(outDir, { recursive: true });
    fs.writeFileSync(outputPath, content, 'utf8');
    await repoGit.add(join(fullPath, outputFile));
  }

  const message = commitMessage;
  await repoGit.commit(message, undefined, {
    '--author': `${commitAuthor.name} <${commitAuthor.email}>`
  });

  await repoGit.push('origin', branch);
}

export async function deleteResource(folderPath) {
  const branch = 'main';
  const commitAuthor = { name: 'Maritime Flow GitOps Bot', email: 'maritimeflowfurb@gmail.com' };
  const commitMessage = `Deleted resource at ${folderPath}`;
  const repoGit = await ensureRepo();
  const resourcePath = join(fullPath, folderPath);

   if (fs.existsSync(resourcePath)) {
    fs.rmSync(resourcePath, { recursive: true, force: true });
  } else {
    // nothing to delete; but we'll still try to ensure git index is clean
  }

  await repoGit.add(['-A', resourcePath]);

  const message = commitMessage;
  await repoGit.commit(message, undefined, {
    '--author': `${commitAuthor.name} <${commitAuthor.email}>`
  });

  await repoGit.push('origin', branch);
}

function escapeRegExp(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
