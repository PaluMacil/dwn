import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { BlogPage, BlogMode } from '../pages/Blog/Blog.page';
import { BlogEditorPage, EditorMode } from '../pages/BlogEditor/BlogEditor.page';

const routes: Routes = [
  {
    path: 'blog/edit/post/:id',
    component: BlogEditorPage,
    data: { mode: EditorMode.Edit }
  },
  {
    path: 'blog/new/post/:topic',
    component: BlogEditorPage,
    data: { mode: EditorMode.New }
  },
  {
    path: 'blog/new/post',
    component: BlogEditorPage,
    data: { mode: EditorMode.New }
  },
  {
    path: 'blog/post/:id/:title',
    component: BlogPage,
    data: { mode: BlogMode.SinglePost }
  },
  {
    path: 'blog/topic/:topicID/:topic',
    component: BlogPage,
    data: { mode: BlogMode.Topic }
  },
  {
    path: '',
    component: BlogPage,
    data: { mode: BlogMode.All }
  }
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes)
  ],
  exports: [RouterModule]
})
export class RoutingModule { }