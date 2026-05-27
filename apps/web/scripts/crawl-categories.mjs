import { writeFileSync } from "fs";
import { join, dirname } from "path";
import { fileURLToPath } from "url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const TIKI_API = "https://tiki.vn/api/v2/categories";

async function fetchJSON(url) {
  const res = await fetch(url);
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${url}`);
  return res.json();
}

function toCategory(tikiCat, depth, children = []) {
  return {
    id: String(tikiCat.id),
    parent_id: tikiCat.parent_id ? String(tikiCat.parent_id) : null,
    name: tikiCat.name,
    slug: tikiCat.url_key,
    description: tikiCat.meta_description || "",
    image_url: tikiCat.thumbnail_url || "",
    sort_order: 1,
    is_active: tikiCat.status === "active",
    depth,
    path: tikiCat.url_path || "",
    product_count: tikiCat.product_count || 0,
    ...(children.length > 0 ? { children } : {}),
  };
}

async function fetchSubCategories(parentId, depth) {
  try {
    const { data } = await fetchJSON(`${TIKI_API}?parent_id=${parentId}`);
    const result = [];
    for (const cat of data) {
      let children = [];
      if (!cat.is_leaf) {
        children = await fetchSubCategories(cat.id, depth + 1);
      }
      result.push(toCategory(cat, depth, children));
    }
    return result;
  } catch {
    return [];
  }
}

async function main() {
  const categoryId = process.argv[2] || "1883";
  console.log(`Crawling Tiki categories (parent_id=${categoryId})...`);

  const root = await fetchJSON(`${TIKI_API}/${categoryId}`);
  const children = await fetchSubCategories(root.id, (root.level || 2) + 1);

  const tree = [toCategory(root, root.level || 2, children)];

  const outPath = join(__dirname, "..", "src", "data", "tiki-categories.json");
  writeFileSync(outPath, JSON.stringify(tree, null, 2), "utf-8");
  console.log(`Written ${outPath}`);
  console.log(`Categories: ${JSON.stringify(tree).match(/"id":"/g)?.length || 0}`);
}

main().catch(console.error);
