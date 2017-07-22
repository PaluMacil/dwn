import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { BlogPage } from '../pages/Blog/Blog.page';

const routes: Routes = [
  {
    path: 'blog/post/:id/:title',
    component: BlogPage
  },
  {
    path: 'blog/topic/:topicID/:topic',
    component: BlogPage
  },
  {
    path: '',
    component: BlogPage
  }
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes)
  ],
  exports: [RouterModule]
})
export class RoutingModule { }