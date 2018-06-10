using Dwn.Data.Models.Identity;
using System;
using System.Collections.Generic;
using System.Text;

namespace Dwn.Data.Models.Blogging
{
    public class BlogCategory
    {
        public Guid Id { get; set; }
        public string Name { get; set; }
        public bool Personal { get; set; }
        public virtual List<Post> Posts { get; set; }
        public virtual List<ApplicationRole> ViewerRoles { get; set; }
    }
}
