using Dwn.Data.Models.Identity;
using System;
using System.Collections.Generic;
using System.Text;

namespace Dwn.Data.Models.Authoring
{
    public class Modification
    {
        public Guid Id { get; set; }

        public ApplicationUser Author { get; set; }
        public DateTime Modified { get; set; }
    }
}